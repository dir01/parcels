package cainiao_test

import (
	"context"
	"encoding/json"
	"github.com/dir01/parcels/externalapis/cainiao"
	"github.com/dir01/parcels/service"
	"os"
	"path"
	"testing"
)

func TestCainiao(t *testing.T) {
	api := cainiao.New()

	t.Run("RS0814398526Y", func(t *testing.T) {
		resp := loadGoldenOrFetch(t, api, "RS0814398526Y")
		t.Logf("response: %+v", resp)
		if resp.Status != service.StatusSuccess {
			t.Fatalf("unexpected status: %v", resp.Status)
		}
		info, err := api.Parse(resp)
		if err != nil {
			t.Fatalf("unexpected error while parsing resp: %v", err)
		}
		if info.TrackingNumber != "RS0814398526Y" {
			t.Fatalf("Unexpected TrackingNumber: %s", info.TrackingNumber)
		}
	})

	t.Run("UZ0556033196Y", func(t *testing.T) {
		resp := loadGoldenOrFetch(t, api, "UZ0556033196Y")
		t.Logf("response: %+v", resp)
		if resp.Status != service.StatusNotFound {
			t.Fatalf("unexpected status: %v", resp.Status)
		}
		info, err := api.Parse(resp)
		if err != nil {
			t.Fatalf("unexpected error while parsing resp: %v", err)
		}
		if info.TrackingNumber != "UZ0556033196Y" {
			t.Fatalf("Unexpected TrackingNumber: %s", info.TrackingNumber)
		}
	})
}

func loadGoldenOrFetch(t *testing.T, api service.PostalAPI, trackingNumber string) service.PostalApiResponse {
	// if UPDATE_TESTDATA in env or file is missing, fetch from API and save to file
	// otherwise, load from file and respond
	// both response and error are supported, so golden file is binary-serialized
	goldenPath := t.Name() + ".golden"

	if info, err := os.Stat(goldenPath); err == nil && info.Size() != 0 && os.Getenv("UPDATE_TESTDATA") == "" {
		bytes, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("failed to read golden file: %v", err)
		}
		var resp service.PostalApiResponse
		if err := json.Unmarshal(bytes, &resp); err != nil {
			t.Fatalf("failed to unmarshal golden file: %v", err)
		}
		return resp
	}

	resp := api.Fetch(context.Background(), trackingNumber)
	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	// mkdir -p $(dirname goldenPath)
	dirname := path.Dir(goldenPath)
	if err := os.MkdirAll(dirname, 0755); err != nil {
		t.Fatalf("failed to create golden file dir:  %v", err)
	}
	if err := os.WriteFile(goldenPath, bytes, 0644); err != nil {
		t.Fatalf("failed to write golden file: %v", err)
	}
	return resp
}
