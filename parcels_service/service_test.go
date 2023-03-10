package parcels_service_test

import (
	"context"
	"testing"
	"time"

	"github.com/dir01/parcels/parcels_service"
	"github.com/dir01/parcels/parcels_service/mocks"
	"go.uber.org/zap"
)

//go:generate minimock -g -i github.com/dir01/parcels/parcels_service.Storage,github.com/dir01/parcels/parcels_service.PostalApi -o ./mocks -s _mock.go
func TestService(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	okCheckInterval := 24 * time.Hour
	notFoundCheckInterval := 3 * 24 * time.Hour
	unknownErrorCheckInterval := 3 * time.Hour
	apiFetchTimeout := 1 * time.Millisecond
	expiryTimeout := 6 * 30 * 24 * time.Hour

	prepareTestSubjects := func() (
		svc *parcels_service.ServiceImpl,
		storage *mocks.StorageMock,
		timeCh chan time.Time,
		api1 *mocks.PostalApiMock,
	) {
		now := time.Now()
		timeCh = make(chan time.Time)
		go func() {
			for newNow := range timeCh {
				now = newNow
			}
		}()

		storage = mocks.NewStorageMock(t)

		api1 = mocks.NewPostalApiMock(t)

		apiMap := map[string]parcels_service.PostalApi{
			"api1": api1,
		}

		svc = parcels_service.NewService(
			apiMap,
			storage,
			okCheckInterval,
			notFoundCheckInterval,
			unknownErrorCheckInterval,
			apiFetchTimeout,
			expiryTimeout,
			logger,
			func() time.Time {
				return now
			},
		)

		return svc, storage, timeCh, api1
	}

	t.Run("new tracking number", func(t *testing.T) {
		callCtx := context.WithValue(context.Background(), "foo", "bar")
		svc, storage, _, api1 := prepareTestSubjects()

		storage.GetLatestMock.Expect(callCtx, "123", []string{"api1"}).Return(nil, nil)

		api1Response := &parcels_service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			Status:         parcels_service.StatusSuccess,
			ResponseBody:   []byte("foo"),
		}
		api1.FetchMock.Inspect(func(ctx context.Context, trackingNumber string) {
			if ctx.Value("foo") != "bar" {
				t.Fatalf("expected context to be inhereted from %v, got %v", callCtx, ctx)
			}
			if trackingNumber != "123" {
				t.Fatalf("expected tracking number to be 123, got %s", trackingNumber)
			}
		}).Return(api1Response)

		storage.InsertMock.Expect(callCtx, "123", "api1", api1Response).Return(nil)

		api1TrackingInfo := &parcels_service.TrackingInfo{
			TrackingNumber: "123",
			ApiName:        "api1",
		}
		api1.ParseMock.Expect(api1Response).Return(api1TrackingInfo, nil)

		tracking, err := svc.GetTrackingInfo(callCtx, "123")
		if err != nil {
			t.Fatalf("failed to get tracking info: %v", err)
		}

		if len(tracking) != 1 {
			t.Fatalf("expected 1 tracking info, got %d", len(tracking))
		}
		if tracking[0] != api1TrackingInfo {
			t.Fatalf("expected tracking info to be %v, got %v", api1TrackingInfo, tracking[0])
		}
	})

	t.Run("freshly fetched tracking number", func(t *testing.T) {
		callCtx := context.WithValue(context.Background(), "foo", "bar")
		svc, storage, timeCh, api1 := prepareTestSubjects()
		defer close(timeCh)
		now := time.Now()
		timeCh <- now
		storedRawResponse := &parcels_service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			Status:         parcels_service.StatusSuccess,
			ResponseBody:   []byte("foo"),
			LastFetchedAt:  now.Add(-(okCheckInterval / 2)),
		}
		storage.GetLatestMock.
			Expect(callCtx, "123", []string{"api1"}).
			Return([]*parcels_service.PostalApiResponse{storedRawResponse}, nil)
		parsedTrackingInfo := &parcels_service.TrackingInfo{
			TrackingNumber: "123",
			ApiName:        "api1",
		}
		api1.ParseMock.Expect(storedRawResponse).Return(parsedTrackingInfo, nil)

		tr, err := svc.GetTrackingInfo(callCtx, "123")
		if err != nil {
			t.Fatalf("failed to get tracking info: %v", err)
		}

		if len(tr) != 1 {
			t.Fatalf("expected 1 tracking info, got %d", len(tr))
		}
		if tr[0] != parsedTrackingInfo {
			t.Fatalf("expected tracking info to be %v, got %v", parsedTrackingInfo, tr[0])
		}
	})
}
