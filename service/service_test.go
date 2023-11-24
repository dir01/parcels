package service_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dir01/parcels/service"
	"github.com/dir01/parcels/service/mocks"
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
		svc *service.ServiceImpl,
		storage *mocks.StorageMock,
		setNow func(t time.Time),
		api1 *mocks.PostalApiMock,
	) {
		now := time.Now()
		setNow = func(t time.Time) {
			now = t
		}
		storage = mocks.NewStorageMock(t)

		api1 = mocks.NewPostalApiMock(t)

		apiMap := map[string]service.PostalApi{
			"api1": api1,
		}

		svc = service.NewService(
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

		return svc, storage, setNow, api1
	}

	t.Run("new tracking number", func(t *testing.T) {
		callCtx := context.WithValue(context.Background(), "foo", "bar")
		svc, storage, setNow, api1 := prepareTestSubjects()

		storage.GetLatestMock.Expect(callCtx, "123", []string{"api1"}).Return(nil, nil)

		api1Response := service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			Status:         service.StatusSuccess,
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

		now := time.Now()
		setNow(now)

		// service should set the timestamps, so from now on we expect them to be set
		api1Response.FirstFetchedAt = now
		api1Response.LastFetchedAt = now

		storage.InsertMock.Expect(callCtx, "123", "api1", &api1Response).Return(nil)

		api1TrackingInfo := &service.TrackingInfo{
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

	t.Run("stored tracking - ttl not expired", func(t *testing.T) {
		callCtx := context.WithValue(context.Background(), "foo", "bar")
		svc, storage, setNow, api1 := prepareTestSubjects()

		now := time.Now()
		setNow(now)

		storedRawResponse := service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			Status:         service.StatusSuccess,
			ResponseBody:   []byte("foo"),
			LastFetchedAt:  now.Add(-(okCheckInterval / 2)),
		}
		storage.GetLatestMock.
			Expect(callCtx, "123", []string{"api1"}).
			Return([]*service.PostalApiResponse{&storedRawResponse}, nil)
		parsedTrackingInfo := &service.TrackingInfo{
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
		expected := &service.TrackingInfo{
			TrackingNumber: "123",
			ApiName:        "api1",
			LastFetchedAt:  now.Add(-(okCheckInterval / 2)),
		}
		if !reflect.DeepEqual(tr[0], expected) {
			t.Fatalf("expected tracking info to be %v, got %v", expected, tr[0])
		}
	})

	t.Run("not found tracking number", func(t *testing.T) {
		callCtx := context.WithValue(context.Background(), "foo", "bar")
		svc, storage, setNow, api1 := prepareTestSubjects()

		storage.GetLatestMock.Expect(callCtx, "123", []string{"api1"}).Return(nil, nil)

		api1Response := service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			Status:         service.StatusNotFound,
			ResponseBody:   []byte("whatever"),
		}
		api1.FetchMock.Inspect(func(ctx context.Context, trackingNumber string) {
			if ctx.Value("foo") != "bar" {
				t.Fatalf("expected context to be inhereted from %v, got %v", callCtx, ctx)
			}
			if trackingNumber != "123" {
				t.Fatalf("expected tracking number to be 123, got %s", trackingNumber)
			}
		}).Return(api1Response)

		now := time.Now()
		setNow(now)

		storage.InsertMock.Expect(callCtx, "123", "api1", &service.PostalApiResponse{
			TrackingNumber: "123",
			ApiName:        "api1",
			FirstFetchedAt: now,
			LastFetchedAt:  now,
			ResponseBody:   []byte("whatever"),
			Status:         service.StatusNotFound,
		}).Return(nil)

		tracking, err := svc.GetTrackingInfo(callCtx, "123")
		if err != nil {
			t.Fatalf("failed to get tracking info: %v", err)
		}

		if len(tracking) != 0 {
			t.Fatalf("expected 0 tracking info, got %d", len(tracking))
		}
	})

}
