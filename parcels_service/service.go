package parcels_service

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

func NewService(
	postalApiMap map[string]PostalApi,
	storage Storage,
	okCheckInterval time.Duration,
	notFoundCheckInterval time.Duration,
	unknownErrorCheckInterval time.Duration,
	apiFetchTimeout time.Duration,
	expiryTimeout time.Duration,
	logger *zap.Logger,
	now func() time.Time,
) *ServiceImpl {
	s := &ServiceImpl{
		apiMap:                    postalApiMap,
		apiNames:                  maps.Keys(postalApiMap),
		storage:                   storage,
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

type Service interface {
	GetTrackingInfo(ctx context.Context, trackingNumber string) ([]*TrackingInfo, error)
}

type ServiceImpl struct {
	apiMap                    map[string]PostalApi
	apiNames                  []string
	storage                   Storage
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
	GetLatest(ctx context.Context, trackingNumber string, apiNames []string) ([]*PostalApiResponse, error)
	Insert(ctx context.Context, trackingNumber string, apiName string, response *PostalApiResponse) error // signature requires trackingNumber and apiName just to add gravity to api contract
	Update(context.Context, *PostalApiResponse) error
}

// PostalApi represents a single postal service API.
// It should know 2 things:
// 1. How to fetch raw response from the postal service's API
// 2. How to parse raw response into TrackingInfo
type PostalApi interface {
	// Fetch should not return error.
	// If there is an error, it should be indicated in PostalApiResponse.Status
	Fetch(ctx context.Context, trackingNumber string) *PostalApiResponse
	Parse(rawResponse *PostalApiResponse) (*TrackingInfo, error)
}

func (svc *ServiceImpl) GetTrackingInfo(ctx context.Context, trackingNumber string) ([]*TrackingInfo, error) {
	now := svc.now()
	storedResponsesMap, err := svc.loadRawResponsesMap(ctx, trackingNumber)
	if err != nil {
		return nil, zaperr.Wrap(err, "loadRawResponsesMap")
	}

	apisToHit, isParcelDelivered, parsedResponsesMap := svc.analyzeStoredResponses(storedResponsesMap)

	if isParcelDelivered {
		return maps.Values(parsedResponsesMap), nil
	}

	fetchedResponsesMap := svc.fetchResponses(ctx, trackingNumber, apisToHit)

	result := make([]*TrackingInfo, 0, len(fetchedResponsesMap)+len(storedResponsesMap))

	// getParsedResp is a convenience function to get parsed response from the map
	// or parse it if it's not there yet
	getParsedResp := func(apiName string, rawResp *PostalApiResponse) *TrackingInfo {
		if parsed, exists := parsedResponsesMap[apiName]; exists {
			return parsed
		} else if parsed, err := svc.parseApiResponse(rawResp); err == nil {
			return parsed
		}
		return nil
	}

	for _, apiName := range svc.apiNames {
		fetched := fetchedResponsesMap[apiName]
		stored := storedResponsesMap[apiName]

		if fetched == nil && stored == nil {
			continue
		}

		switch {
		case fetched == nil:
			// Stored response, was too fresh to fetch, no need to update, just return it
			if parsed := getParsedResp(apiName, stored); parsed != nil {
				result = append(result, parsed)
			}
		case stored == nil || !bytes.Equal(stored.ResponseBody, fetched.ResponseBody):
			// Something new, store it and return it
			fetched.FirstFetchedAt = now
			fetched.LastFetchedAt = now
			if err := svc.storage.Insert(ctx, trackingNumber, apiName, fetched); err != nil {
				// TODO: check for duplicate key error, find correct entry and update it
			}
			if parsed := getParsedResp(apiName, fetched); parsed != nil {
				result = append(result, parsed)
			}
		default: // stored.ResponseBody == fetched.ResponseBody
			// We already have this response, just update the timestamp
			stored.LastFetchedAt = now
			if err := svc.storage.Update(ctx, stored); err != nil {
				svc.log.Error("failed to update stored response", zap.Error(err))
			}
			if parsed := getParsedResp(apiName, stored); parsed != nil {
				result = append(result, parsed)
			}
		}
	}

	return result, nil
}

// loadRawResponsesMap loads the last responses from all APIs,
// and returns them as a map[apiName]response
func (svc *ServiceImpl) loadRawResponsesMap(ctx context.Context, trackingNumber string) (map[string]*PostalApiResponse, error) {
	lastResponses, err := svc.storage.GetLatest(ctx, trackingNumber, svc.apiNames)
	if err != nil {
		return nil, zaperr.Wrap(err, "failed to get latest responses")
	}
	lastResponsesMap := make(map[string]*PostalApiResponse, len(lastResponses))
	for _, resp := range lastResponses {
		lastResponsesMap[resp.ApiName] = resp
	}
	return lastResponsesMap, nil
}

// analyzeStoredResponses goes through the last responses from all APIs
// and returns:
// - apisToHit: list of APIs that should be re-fetched for new responses
// - isDelivered: whether the package is delivered
// - parsedResponsesMap: result of responses parsing
func (svc *ServiceImpl) analyzeStoredResponses(
	lastRespMap map[string]*PostalApiResponse,
) (
	apisToHit []string,
	isParcelDelivered bool,
	parsedResponsesMap map[string]*TrackingInfo,
) {
	parsedResponsesMap = make(map[string]*TrackingInfo, len(lastRespMap))
	isParcelDelivered = false
	apiHitDecisionMap := make(map[string]bool, len(svc.apiNames))
	for _, apiName := range svc.apiNames {
		apiHitDecisionMap[apiName] = false
		resp := lastRespMap[apiName]

		if resp == nil {
			// new tracking numbers we've never seen before
			apiHitDecisionMap[apiName] = true
			continue
		}

		if parsed, err := svc.parseApiResponse(resp); err == nil {
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
			isCheckTimeoutPassed := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = isCheckTimeoutPassed
			continue
		case resp.Status == StatusUnknownError:
			// Last time we got an error from this API.
			// This tells us nothing about relevance of this API.
			// Should be checked less often.
			recheckAt := resp.LastFetchedAt.Add(svc.unknownErrorCheckInterval)
			isCheckTimeoutPassed := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = isCheckTimeoutPassed
		case resp.Status == StatusNotFound:
			// API never indicated that it knows about this tracking number.
			// Most likely, it's not a relevant API.
			// Should be checked least often.
			recheckAt := resp.LastFetchedAt.Add(svc.notFoundCheckInterval)
			isCheckTimeoutPassed := svc.now().After(recheckAt)
			apiHitDecisionMap[apiName] = isCheckTimeoutPassed
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
func (svc *ServiceImpl) fetchResponses(ctx context.Context, trackingNumber string, apisToHit []string) map[string]*PostalApiResponse {
	fetchedResponsesMap := make(map[string]*PostalApiResponse, len(apisToHit))

	resultsChan := make(chan *PostalApiResponse, len(apisToHit))
	wg := sync.WaitGroup{}

	ttlCtx, cancel := context.WithTimeout(ctx, svc.apiFetchTimeout)
	defer cancel()
	for _, apiName := range apisToHit {
		wg.Add(1)
		go func(apiName string) {
			defer wg.Done()
			resp := svc.apiMap[apiName].Fetch(ttlCtx, trackingNumber)
			resultsChan <- resp
		}(apiName)
	}

	wg.Wait()
	for i := 0; i < len(apisToHit); i++ {
		var resp *PostalApiResponse
		select {
		case <-ttlCtx.Done(): // yes, we pass the context to the API,
			// but can we really trust them to respect it?
			return fetchedResponsesMap
		case resp = <-resultsChan:
			break
		}

		if resp == nil {
			svc.log.Debug(
				"nil response from API",
				zap.String("api", resp.ApiName),
				zap.String("trackingNumber", trackingNumber),
			)
			continue
		}
		fetchedResponsesMap[resp.ApiName] = resp
	}
	close(resultsChan)

	return fetchedResponsesMap
}

func (svc *ServiceImpl) parseApiResponse(resp *PostalApiResponse) (*TrackingInfo, error) {
	if resp.Status != StatusSuccess {
		return nil, fmt.Errorf("not parsing a response with a non-success status")
	}
	return svc.apiMap[resp.ApiName].Parse(resp)
}
