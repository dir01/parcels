package service

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hori-ryota/zaperr"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type APIName string

type Service interface {
	GetTrackingInfo(ctx context.Context, trackingNumber string) ([]*TrackingInfo, error)
}

func NewService(
	postalApiMap map[APIName]PostalAPI,
	storage Storage,
	metrics Metrics,
	okCheckInterval time.Duration,
	notFoundCheckInterval time.Duration,
	unknownErrorCheckInterval time.Duration,
	apiFetchTimeout time.Duration,
	expiryTimeout time.Duration,
	logger *zap.Logger,
	now func() time.Time,
) *Impl {
	s := &Impl{
		apiMap:                    postalApiMap,
		apiNames:                  maps.Keys(postalApiMap),
		storage:                   storage,
		metrics:                   metrics,
		okCheckInterval:           okCheckInterval,
		notFoundCheckInterval:     notFoundCheckInterval,
		unknownErrorCheckInterval: unknownErrorCheckInterval,
		apiFetchTimeout:           apiFetchTimeout,
		expiryTimeout:             expiryTimeout,
		log:                       logger,
		now:                       now,
	}
	var _ Service = s
	return s
}

type Impl struct {
	apiMap                    map[APIName]PostalAPI
	apiNames                  []APIName
	storage                   Storage
	metrics                   Metrics
	okCheckInterval           time.Duration
	notFoundCheckInterval     time.Duration
	unknownErrorCheckInterval time.Duration
	expiryTimeout             time.Duration
	log                       *zap.Logger
	now                       func() time.Time
	apiFetchTimeout           time.Duration
}

// Storage contains whole history of PostalAPI responses.
// It is used to avoid unnecessary calls to the API.
// We store whole history of responses for further analysis.
// However, this storage is only concerned with the last response.
type Storage interface {
	GetLatest(ctx context.Context, trackingNumber string, apiNames []APIName) ([]*PostalApiResponse, error)
	// Insert signature requires trackingNumber and apiName just to add gravity to api contract
	// PostalApiResponse could have no
	Insert(ctx context.Context, trackingNumber string, apiName APIName, response *PostalApiResponse) error
	Update(context.Context, *PostalApiResponse) error
}

// Metrics describes what custom metrics service should report on
type Metrics interface {
	ParcelDelivered()

	FetchedChanged(APIName)
	FetchedFirst(APIName)
	FetchedUnchanged(APIName)

	APIHit(APIName)
	APIParseError(APIName)

	CacheBustAfterSuccess(apiName APIName, willRefetch bool)
	CacheBustAfterUnknownError(apiName APIName, willRefetch bool)
	CacheBustAfterNotFoundError(apiName APIName, willRefetch bool)
}

// PostalAPI represents a single postal service API.
// It should know 2 things:
// 1. How to fetch raw response from the postal service's API
// 2. How to parse raw response into TrackingInfo
type PostalAPI interface {
	// Fetch should not return error, and response should never be nil.
	// Any error and missing response should be indicated in PostalApiResponse.Status
	Fetch(ctx context.Context, trackingNumber string) PostalApiResponse
	Parse(rawResponse PostalApiResponse) (*TrackingInfo, error)
}

func (svc *Impl) GetTrackingInfo(ctx context.Context, trackingNumber string) ([]*TrackingInfo, error) {
	now := svc.now()
	storedResponsesMap, err := svc.loadRawResponsesMap(ctx, trackingNumber)
	svc.log.Info(
		"loaded stored responses",
		zap.String("trackingNumber", trackingNumber),
		zap.Int("count", len(storedResponsesMap)),
	)
	if err != nil {
		return nil, zaperr.Wrap(err, "loadRawResponsesMap")
	}

	apisToHit, isParcelDelivered, parsedResponsesMap := svc.analyzeStoredResponses(storedResponsesMap)

	if isParcelDelivered {
		svc.log.Info(
			"parcel is delivered",
			zap.String("trackingNumber", trackingNumber),
		)
		svc.metrics.ParcelDelivered()
		return maps.Values(parsedResponsesMap), nil
	}

	fetchedResponsesMap := svc.fetchResponses(ctx, trackingNumber, apisToHit)

	result := make([]*TrackingInfo, 0, len(fetchedResponsesMap)+len(storedResponsesMap))

	// getParsedResp is a convenience function to get parsed response from the map
	// or parse it if it's not there yet
	getParsedResp := func(apiName APIName, rawResp PostalApiResponse) *TrackingInfo {
		if parsed, exists := parsedResponsesMap[apiName]; exists {
			return parsed
		} else if parsed, err := svc.parseApiResponse(rawResp); err == nil {
			return parsed
		} else if err != nil {
			svc.metrics.APIParseError(apiName)
		}
		return nil
	}

	for _, apiName := range svc.apiNames {
		fetched, wasFetched := fetchedResponsesMap[apiName]
		stored := storedResponsesMap[apiName]

		switch {
		case stored != nil && !wasFetched:
			// Stored response was found, but was too fresh to re-fetch, no need to update, just return it
			if parsed := getParsedResp(apiName, *stored); parsed != nil { // we can dereference here because we know that stored != nil
				parsed.LastFetchedAt = stored.LastFetchedAt
				result = append(result, parsed)
			}
		case stored == nil || !bytes.Equal(stored.ResponseBody, fetched.ResponseBody):
			if stored != nil {
				svc.metrics.FetchedChanged(apiName)
			} else {
				svc.metrics.FetchedFirst(apiName)
			}
			// Something new, store it and return it (but only if it's not an error)
			fetched.FirstFetchedAt = now
			fetched.LastFetchedAt = now
			if err := svc.storage.Insert(ctx, trackingNumber, apiName, &fetched); err != nil {
				// TODO: check for duplicate key error, find correct entry and update it
			}
			if fetched.Status == StatusSuccess {
				if parsed := getParsedResp(apiName, fetched); parsed != nil {
					result = append(result, parsed)
				}
			}
		default: // stored != nil && stored.ResponseBody == fetched.ResponseBody
			// We already have this response, just update the timestamp
			svc.metrics.FetchedUnchanged(apiName)
			stored.LastFetchedAt = now
			if err := svc.storage.Update(ctx, stored); err != nil {
				svc.log.Error("failed to update stored response", zap.Error(err))
			}
			if parsed := getParsedResp(apiName, *stored); parsed != nil { // we can dereference here because we know that stored != nil
				result = append(result, parsed)
			}
		}
	}

	return result, nil
}

// loadRawResponsesMap loads the last responses from all APIs,
// and returns them as a map[apiName]response
func (svc *Impl) loadRawResponsesMap(
	ctx context.Context,
	trackingNumber string,
) (map[APIName]*PostalApiResponse, error) {
	lastResponses, err := svc.storage.GetLatest(ctx, trackingNumber, svc.apiNames)
	if err != nil {
		return nil, zaperr.Wrap(err, "failed to get latest responses")
	}
	lastResponsesMap := make(map[APIName]*PostalApiResponse, len(lastResponses))
	for _, resp := range lastResponses {
		lastResponsesMap[resp.APIName] = resp
	}
	return lastResponsesMap, nil
}

// analyzeStoredResponses goes through the last responses from all APIs
// and returns:
// - apisToHit: list of APIs that should be re-fetched for new responses
// - isDelivered: whether the package is delivered
// - parsedResponsesMap: result of responses parsing
func (svc *Impl) analyzeStoredResponses(
	lastRespMap map[APIName]*PostalApiResponse,
) (
	apisToHit []APIName,
	isParcelDelivered bool,
	parsedResponsesMap map[APIName]*TrackingInfo,
) {
	parsedResponsesMap = make(map[APIName]*TrackingInfo, len(lastRespMap))
	isParcelDelivered = false
	apiHitDecisionMap := make(map[APIName]bool, len(svc.apiNames))
	for _, apiName := range svc.apiNames {
		apiHitDecisionMap[apiName] = false
		resp := lastRespMap[apiName]

		if resp == nil {
			// new tracking numbers we've never seen before
			apiHitDecisionMap[apiName] = true
			continue
		}

		if parsed, err := svc.parseApiResponse(*resp); err == nil { // we can dereference here because we know that resp != nil
			parsedResponsesMap[apiName] = parsed
			if !isParcelDelivered && parsed.IsDelivered() {
				isParcelDelivered = true
			}
		}

		if isParcelDelivered {
			continue
			// We don't need to analyze whether to hit the APIs or not: we won't.
			// But we do need to continue parsing the responses
			// So we can't break, but we can continue
		}

		switch {
		case svc.now().After(resp.LastFetchedAt.Add(svc.expiryTimeout)): // identical to `nil`
			// Tracking number we've seen long time ago.
			// This can be a tracking number reuse (it happens),
			// so we should treat is as a new tracking number
			apiHitDecisionMap[apiName] = true
			continue
		case resp.Status == StatusSuccess:
			// We already got updates for this tracking number.
			// So we know that this API is a relevant one.
			// Should be checked most often.
			recheckAt := resp.LastFetchedAt.Add(svc.okCheckInterval)
			shouldRefetch := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = shouldRefetch
			svc.metrics.CacheBustAfterSuccess(apiName, shouldRefetch)
			continue
		case resp.Status == StatusUnknownError:
			// Last time we got an error from this API.
			// This tells us nothing about relevance of this API.
			// Should be checked less often.
			recheckAt := resp.LastFetchedAt.Add(svc.unknownErrorCheckInterval)
			shouldRefetch := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = shouldRefetch
			svc.metrics.CacheBustAfterUnknownError(apiName, shouldRefetch)
		case resp.Status == StatusNotFound:
			// API never indicated that it knows about this tracking number.
			// Most likely, it's not a relevant API.
			// Should be checked least often.
			recheckAt := resp.LastFetchedAt.Add(svc.notFoundCheckInterval)
			shouldRefetch := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = shouldRefetch
			svc.metrics.CacheBustAfterNotFoundError(apiName, shouldRefetch)
		default:
			// TODO: check for mentioned countries, and reconsider relevance of the API
		}
	}

	if isParcelDelivered {
		return nil, true, parsedResponsesMap
	}

	for apiName, shouldHit := range apiHitDecisionMap {
		if shouldHit {
			apisToHit = append(apisToHit, apiName)
		}
	}

	return apisToHit, isParcelDelivered, parsedResponsesMap
}

// fetchResponses fetches responses from all the APIs in parallel.
// It does not return any errors; any errors are reflected in responses.
func (svc *Impl) fetchResponses(
	ctx context.Context,
	trackingNumber string,
	apisToHit []APIName,
) map[APIName]PostalApiResponse {
	fetchedResponsesMap := make(map[APIName]PostalApiResponse, len(apisToHit))

	resultsChan := make(chan PostalApiResponse, len(apisToHit))
	wg := sync.WaitGroup{}

	ttlCtx, cancel := context.WithTimeout(ctx, svc.apiFetchTimeout)
	defer cancel()
	for _, apiName := range apisToHit {
		wg.Add(1)
		go func(apiName APIName) {
			defer wg.Done()

			svc.metrics.APIHit(apiName)
			resp := svc.apiMap[apiName].Fetch(ttlCtx, trackingNumber)
			resultsChan <- resp
		}(apiName)
	}

	wg.Wait()
	for i := 0; i < len(apisToHit); i++ {
		select {
		case <-ttlCtx.Done(): // yes, we pass the context to the API,
			// but can we really trust them to respect it?
			return fetchedResponsesMap
		case resp := <-resultsChan:
			fetchedResponsesMap[resp.APIName] = resp
			continue
		}
	}
	close(resultsChan)

	return fetchedResponsesMap
}

func (svc *Impl) parseApiResponse(resp PostalApiResponse) (*TrackingInfo, error) {
	if resp.Status != StatusSuccess {
		return nil, fmt.Errorf("not parsing a response with a non-success status")
	}
	return svc.apiMap[resp.APIName].Parse(resp)
}
