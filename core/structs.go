package core

import (
	"time"
)

// TrackingInfo represents final result of the service:
// parsed and normalized representation of parcel tracking info.
type TrackingInfo struct {
	TrackingNumber            string
	ApiName                   string
	LastFetchedAt             time.Time
	OriginCountry             string
	DestinationCountry        string
	Events                    []TrackingEvent
	AdditionalTrackingNumbers []string
}

func (ti *TrackingInfo) IsDelivered() bool {
	// iterating backwards, since, `DELIVERED` event is most likely the last one
	for i := len(ti.Events) - 1; i >= 0; i-- {
		if ti.Events[i].Status == TrackingStatusDelivered {
			return true
		}
	}
	return false
}

type TrackingEvent struct {
	Time        time.Time
	Description string
	Status      TrackingStatus
}

type TrackingStatus string

const (
	TrackingStatusShipmentInfoReceived          TrackingStatus = "SHIPMENT_INFO_RECEIVED"
	TrackingStatusPackagingComplete             TrackingStatus = "PACKAGING_COMPLETE"
	TrackingStatusDispatchedFromWarehouse       TrackingStatus = "DISPATCHED_FROM_WAREHOUSE"
	TrackingStatusWMSConfirmed                  TrackingStatus = "WMS_CONFIRMED" // WMS = Warehouse Management System
	TrackingStatusArrivedAtSortingCenter        TrackingStatus = "ARRIVED_AT_SORTING_CENTER"
	TrackingStatusAcceptedByCarrier             TrackingStatus = "ACCEPTED_BY_CARRIER"
	TrackingStatusDepartedFromSortingCenter     TrackingStatus = "DEPARTED_FROM_SORTING_CENTER"
	TrackingStatusArrivedAtDepartureHub         TrackingStatus = "ARRIVED_AT_DEPARTURE_TRANSPORT_HUB"
	TrackingStatusTransitPortRerouteCb          TrackingStatus = "TRANSIT_PORT_REROUTE_CALLBACK" //preMainCode:SINOA00241668IL TODO: What does this mean?
	TrackingStatusExportCustomsClearanceStarted TrackingStatus = "EXPORT_CUSTOMS_CLEARANCE_STARTED"
	TrackingStatusLeavignDepartureRegion        TrackingStatus = "LEAVING_FROM_DEPARTURE_COUNTRY_OR_REGION"
	TrackingStatusImportCustomsClearanceStarted TrackingStatus = "IMPORT_CUSTOMS_CLEARANCE_STARTED"
	TrackingStatusImportCustomsClearanceSuccess TrackingStatus = "IMPORT_CUSTOMS_CLEARANCE_SUCCESS"
	TrackingStatusDepartedOriginRegion          TrackingStatus = "DEPARTED_ORIGIN_COUNTRY_OR_REGION"
	TrackingStatusArrivedAtLinehaulOffice       TrackingStatus = "ARRIVED_AT_LINEHAUL_OFFICE"
	TrackingStatusArrivedAtCustoms              TrackingStatus = "ARRIVED_AT_CUSTOMS"
	TrackingStatusDepartedFromCustoms           TrackingStatus = "DEPARTED_FROM_CUSTOMS"
	TrackingStatusExportCustomsClearanceSuccess TrackingStatus = "EXPORT_CUSTOMS_CLEARANCE_SUCCESS"
	TrackingStatusDelivered                     TrackingStatus = "DELIVERED"
	TrackingStatusUnknown                       TrackingStatus = "UNKNOWN"
)

// PostalApiResponse represents raw response from the API.
// We could have used raw strings and known exceptions,
// but I wanted to make API implementations contractually obliged
// to indicate the status of the response while exposing possible expected error types
type PostalApiResponse struct {
	ID             int64
	TrackingNumber string
	ApiName        string
	FirstFetchedAt time.Time
	LastFetchedAt  time.Time
	ResponseBody   []byte
	Status         ApiResponseStatus
}

type ApiResponseStatus string

const (
	StatusSuccess           ApiResponseStatus = "success"
	StatusRateLimitExceeded ApiResponseStatus = "rate_limit_exceeded"
	StatusNotFound          ApiResponseStatus = "not_found"
	StatusUnknownError      ApiResponseStatus = "unknown_error"
)
