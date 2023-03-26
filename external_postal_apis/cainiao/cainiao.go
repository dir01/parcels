package cainiao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dir01/parcels/parcels_service"
)

func New() parcels_service.PostalApi {
	return &Cainiao{}
}

type Cainiao struct{}

func (c *Cainiao) Fetch(ctx context.Context, trackingNumber string) parcels_service.PostalApiResponse {
	result := parcels_service.PostalApiResponse{
		TrackingNumber: trackingNumber,
		ApiName:        "cainiao",
	}

	url := fmt.Sprintf("https://global.cainiao.com/global/detail.json?mailNos=%s&lang=en-US", trackingNumber)
	resp, err := http.Get(url)
	if err != nil {
		result.Status = parcels_service.StatusUnknownError
		return result
	}

	if resp.StatusCode == http.StatusNotFound {
		result.Status = parcels_service.StatusNotFound
		return result
	}

	if resp.StatusCode != http.StatusOK {
		result.Status = parcels_service.StatusUnknownError
		if bytes, err := io.ReadAll(resp.Body); err == nil {
			result.ResponseBody = bytes
		}
		return result
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = parcels_service.StatusUnknownError
		return result
	}

	result.ResponseBody = responseBody
	result.Status = parcels_service.StatusSuccess

	return result
}

func (c *Cainiao) Parse(rawResponse parcels_service.PostalApiResponse) (*parcels_service.TrackingInfo, error) {
	var cainiaoResponse response
	if err := json.Unmarshal(rawResponse.ResponseBody, &cainiaoResponse); err != nil {
		return nil, err
	}

	if len(cainiaoResponse.Module) != 1 {
		return nil, fmt.Errorf("response module length is not 1")
	}

	m0 := cainiaoResponse.Module[0]
	var events []parcels_service.TrackingEvent
	for _, detail := range m0.DetailList {
		if trackingEvent := c.parseDetail(detail); trackingEvent != nil {
			events = append(events, *trackingEvent)
		}
	}

	return &parcels_service.TrackingInfo{
		TrackingNumber:     rawResponse.TrackingNumber,
		ApiName:            rawResponse.ApiName,
		OriginCountry:      m0.OriginCountry,
		DestinationCountry: m0.DestCountry,
		Events:             events,
	}, nil
}

func (c *Cainiao) parseDetail(detail detail) *parcels_service.TrackingEvent {
	return &parcels_service.TrackingEvent{
		Time:        time.Unix(detail.Time/1000, 0),
		Description: detail.StanderdDesc,
		Status:      c.mapStatus(detail.ActionCode),
	}
}

func (c *Cainiao) mapStatus(actionCode string) parcels_service.TrackingStatus {
	switch actionCode {
	case "GWMS_ACCEPT":
		return parcels_service.TrackingStatusShipmentInfoReceived
	case "GWMS_PACKAGE":
		return parcels_service.TrackingStatusPackagingComplete
	case "GWMS_OUTBOUND":
		return parcels_service.TrackingStatusDispatchedFromWarehouse
	case "WMS_CONFIRMED":
		return parcels_service.TrackingStatusWMSConfirmed
	case "SC_INBOUND_SUCCESS":
		return parcels_service.TrackingStatusArrivedAtSortingCenter
	case "PU_PICKUP_SUCCESS":
		return parcels_service.TrackingStatusAcceptedByCarrier
	case "SC_OUTBOUND_SUCCESS":
		return parcels_service.TrackingStatusDepartedFromSortingCenter
	case "LH_HO_IN_SUCCESS":
		return parcels_service.TrackingStatusArrivedAtDepartureHub
	case "TRANSIT_PORT_REROUTE_CALLBACK":
		return parcels_service.TrackingStatusTransitPortRerouteCb
	case "CC_EX_START":
		return parcels_service.TrackingStatusExportCustomsClearanceStarted
	case "CC_EX_SUCCESS":
		return parcels_service.TrackingStatusExportCustomsClearanceSuccess
	case "LH_HO_AIRLINE":
		return parcels_service.TrackingStatusLeavignDepartureRegion
	case "CC_IM_START":
		return parcels_service.TrackingStatusImportCustomsClearanceStarted
	case "LH_DEPART":
		return parcels_service.TrackingStatusDepartedOriginRegion
	case "LH_ARRIVE":
		return parcels_service.TrackingStatusArrivedAtLinehaulOffice
	case "CC_HO_IN_SUCCESS":
		return parcels_service.TrackingStatusArrivedAtCustoms
	case "CC_HO_OUT_SUCCESS":
		return parcels_service.TrackingStatusDepartedFromCustoms
	case "CC_IM_SUCCESS":
		return parcels_service.TrackingStatusImportCustomsClearanceSuccess
	case "CUSTOMS_ARRIVED_IN_AREA_CALLBACK":
		return parcels_service.TrackingStatusArrivedAtCustoms
	default:
		return parcels_service.TrackingStatusUnknown
	}
}
