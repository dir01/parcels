package sqlite_storage

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dir01/parcels/service"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

func TestStorage(t *testing.T) {
	prepareTestSubject := func() service.Storage {
		db := sqlx.MustConnect("sqlite3", ":memory:")
		storage := NewStorage(db)
		migrations := &migrate.FileMigrationSource{
			Dir: "../db/migrations",
		}
		_, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
		if err != nil {
			t.Fatalf("failed to apply migrations: %v", err)
		}
		return storage
	}

	t.Run("Insert and GetLatest", func(t *testing.T) {
		storage := prepareTestSubject()

		rawResp := &service.PostalApiResponse{
			ID:             0,
			TrackingNumber: "some-tracking-number",
			APIName:        "some-api-name",
			FirstFetchedAt: time.Unix(1000, 0),
			LastFetchedAt:  time.Unix(2000, 0),
			ResponseBody:   []byte("some-response-body"),
			Status:         service.StatusSuccess,
		}

		if err := storage.Insert(context.TODO(), "some-tracking-number", "some-api-name", rawResp); err != nil {
			t.Fatalf("failed to insert: %v", err)
		}

		latest, err := storage.GetLatest(context.TODO(), "some-tracking-number", []service.APIName{"some-api-name"})
		if err != nil {
			t.Fatalf("failed to get latest: %v", err)
		}
		if latest == nil {
			t.Fatalf("expected response to be non-nil")
		}
		if len(latest) != 1 {
			t.Fatalf("expected 1 response, got %d", len(latest))
		}
		fetchedResp := latest[0]
		if fetchedResp.TrackingNumber != rawResp.TrackingNumber {
			t.Fatalf("expected tracking number to be %s, got %s", rawResp.TrackingNumber, fetchedResp.TrackingNumber)
		}
		if fetchedResp.APIName != rawResp.APIName {
			t.Fatalf("expected api name to be %s, got %s", rawResp.APIName, fetchedResp.APIName)
		}
		if fetchedResp.FirstFetchedAt != rawResp.FirstFetchedAt {
			t.Fatalf("expected first fetched at to be %s, got %s", rawResp.FirstFetchedAt, fetchedResp.FirstFetchedAt)
		}
		if fetchedResp.LastFetchedAt != rawResp.LastFetchedAt {
			t.Fatalf("expected last fetched at to be %s, got %s", rawResp.LastFetchedAt, fetchedResp.LastFetchedAt)
		}
		if !bytes.Equal(fetchedResp.ResponseBody, rawResp.ResponseBody) {
			t.Fatalf("expected response body to be %s, got %s", rawResp.ResponseBody, fetchedResp.ResponseBody)
		}
		if fetchedResp.Status != rawResp.Status {
			t.Fatalf("expected status to be %s, got %s", rawResp.Status, fetchedResp.Status)
		}
	})

	t.Run("Insert, Update and GetLatest", func(t *testing.T) {
		storage := prepareTestSubject()

		rawResp := &service.PostalApiResponse{
			ID:             1,
			TrackingNumber: "some-tracking-number",
			APIName:        "some-api-name",
			FirstFetchedAt: time.Unix(1000, 0),
			LastFetchedAt:  time.Unix(2000, 0),
			ResponseBody:   []byte("some-response-body"),
			Status:         service.StatusSuccess,
		}

		if err := storage.Insert(context.TODO(), "some-tracking-number", "some-api-name", rawResp); err != nil {
			t.Fatalf("failed to insert: %v", err)
		}

		rawResp.FirstFetchedAt = time.Unix(2000, 0)
		rawResp.LastFetchedAt = time.Unix(3000, 0)
		rawResp.ResponseBody = []byte("some-updated-response-body")
		rawResp.Status = service.StatusRateLimitExceeded

		if err := storage.Update(context.TODO(), rawResp); err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		latest, err := storage.GetLatest(context.TODO(), "some-tracking-number", []service.APIName{"some-api-name"})
		if err != nil {
			t.Fatalf("failed to get latest: %v", err)
		}
		if latest == nil {
			t.Fatalf("expected response to be non-nil")
		}
		if len(latest) != 1 {
			t.Fatalf("expected 1 response, got %d", len(latest))
		}
		fetchedResp := latest[0]
		if fetchedResp.TrackingNumber != rawResp.TrackingNumber {
			t.Fatalf("expected tracking number to be %s, got %s", rawResp.TrackingNumber, fetchedResp.TrackingNumber)
		}
		if fetchedResp.APIName != rawResp.APIName {
			t.Fatalf("expected api name to be %s, got %s", rawResp.APIName, fetchedResp.APIName)
		}
		if fetchedResp.FirstFetchedAt != rawResp.FirstFetchedAt {
			t.Fatalf("expected first fetched at to be %s, got %s", rawResp.FirstFetchedAt, fetchedResp.FirstFetchedAt)
		}
		if fetchedResp.LastFetchedAt != rawResp.LastFetchedAt {
			t.Fatalf("expected last fetched at to be %s, got %s", rawResp.LastFetchedAt, fetchedResp.LastFetchedAt)
		}
		if !bytes.Equal(fetchedResp.ResponseBody, rawResp.ResponseBody) {
			t.Fatalf("expected response body to be %s, got %s", rawResp.ResponseBody, fetchedResp.ResponseBody)
		}
		if fetchedResp.Status != rawResp.Status {
			t.Fatalf("expected status to be %s, got %s", rawResp.Status, fetchedResp.Status)
		}

	})

	t.Run("Insert respects context", func(t *testing.T) {

		storage := prepareTestSubject()
		ttlCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		rawResp := &service.PostalApiResponse{ID: 1}

		err := storage.Insert(ttlCtx, "some-tracking-number", "some-api-name", rawResp)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context deadline exceeded, got %v", err)
		}
	})

	t.Run("GetLatest respects context", func(t *testing.T) {
		storage := prepareTestSubject()
		ttlCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		_, err := storage.GetLatest(ttlCtx, "some-tracking-number", []service.APIName{"some-api-name"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context deadline exceeded, got %v", err)
		}
	})

	t.Run("Update respects context", func(t *testing.T) {
		storage := prepareTestSubject()
		ttlCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		rawResp := &service.PostalApiResponse{ID: 1}

		err := storage.Update(ttlCtx, rawResp)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected context deadline exceeded, got %v", err)
		}
	})
}
