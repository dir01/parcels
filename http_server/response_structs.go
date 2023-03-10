package http_server

import (
	"time"

	"github.com/dir01/parcels/core"
)

type HttpTrackingInfo struct {
	TrackingNumber string `json:"tracking_number"`
	ApiName        string `json:"api_name"`
	IsDelivered    bool   `json:"is_delivered"`
	LastCheckedAt  string `json:"last_checked_at"`
	TrackingEvents []HttpTrackingEvent
}

type HttpTrackingEvent struct {
	Time        string `json:"time"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (hti HttpTrackingInfo) FromBusinessStruct(t *core.TrackingInfo) *HttpTrackingInfo {
	hti.TrackingNumber = t.TrackingNumber
	hti.ApiName = t.ApiName
	hti.IsDelivered = t.IsDelivered()
	hti.LastCheckedAt = t.LastFetchedAt.Format(time.RFC3339)
	for _, e := range t.Events {
		hti.TrackingEvents = append(hti.TrackingEvents, HttpTrackingEvent{
			Time:        e.Time.Format(time.RFC3339),
			Description: e.Description,
			Status:      string(e.Status),
		})
	}
	return &hti
}
