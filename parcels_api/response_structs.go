package parcels_api

import (
	"time"

	"github.com/dir01/parcels/parcels_service"
)

// TrackingInfo represents a single track of a parcel according to one carrier in an API response
type TrackingInfo struct {
	TrackingNumber string          `json:"tracking_number"`
	ApiName        string          `json:"api_name"`
	IsDelivered    bool            `json:"is_delivered"`
	LastCheckedAt  string          `json:"last_checked_at"`
	LastUpdatedAt  string          `json:"last_updated_at"`
	Events         []TrackingEvent `json:"events"`
}

// TrackingEvent represents a single event in a parcel's track
type TrackingEvent struct {
	Time        string `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (hti TrackingInfo) fromBusinessStruct(t *parcels_service.TrackingInfo) *TrackingInfo {
	hti.TrackingNumber = t.TrackingNumber
	hti.ApiName = t.ApiName
	hti.IsDelivered = t.IsDelivered()
	hti.LastCheckedAt = t.LastFetchedAt.Format(time.RFC3339)
	maxTime := time.Time{}
	for _, e := range t.Events {
		if e.Time.After(maxTime) {
			maxTime = e.Time
		}
		hti.Events = append(hti.Events, TrackingEvent{
			Time:        e.Time.Format(time.RFC3339),
			Description: e.Description,
			Status:      string(e.Status),
		})
	}
	hti.LastUpdatedAt = maxTime.Format(time.RFC3339)
	return &hti
}
