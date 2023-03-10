package cainiao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dir01/parcels/core"
)

func New() core.PostalApi {
	return &Cainiao{}
}

type Cainiao struct{}

func (c *Cainiao) Fetch(ctx context.Context, trackingNumber string) *core.PostalApiResponse {
	result := &core.PostalApiResponse{
		TrackingNumber: trackingNumber,
		ApiName:        "cainiao",
	}

	url := fmt.Sprintf("https://global.cainiao.com/global/detail.json?mailNos=%s&lang=en-US", trackingNumber)
	resp, err := http.Get(url)
	if err != nil {
		result.Status = core.StatusUnknownError
		return result
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = core.StatusUnknownError
		return result
	}

	result.ResponseBody = responseBody
	result.Status = core.StatusSuccess

	return result
}

func (c *Cainiao) Parse(rawResponse *core.PostalApiResponse) (*core.TrackingInfo, error) {
	var cainiaoResponse response
	if err := json.Unmarshal(rawResponse.ResponseBody, &cainiaoResponse); err != nil {
		return nil, err
	}

	if len(cainiaoResponse.Module) != 1 {
		return nil, fmt.Errorf("response module length is not 1")
	}

	m0 := cainiaoResponse.Module[0]
	var events []core.TrackingEvent
	for _, detail := range m0.DetailList {
		if trackingEvent := c.parseDetail(detail); trackingEvent != nil {
			events = append(events, *trackingEvent)
		}
	}

	return &core.TrackingInfo{
		TrackingNumber:     rawResponse.TrackingNumber,
		ApiName:            rawResponse.ApiName,
		OriginCountry:      m0.OriginCountry,
		DestinationCountry: m0.DestCountry,
		Events:             events,
	}, nil
}

func (c *Cainiao) parseDetail(detail detail) *core.TrackingEvent {
	return &core.TrackingEvent{
		Time:        time.Unix(detail.Time/1000, 0),
		Description: detail.StanderdDesc,
		Status:      c.mapStatus(detail.ActionCode),
	}
}

func (c *Cainiao) mapStatus(actionCode string) core.TrackingStatus {
	switch actionCode {
	case "GWMS_ACCEPT":
		return core.TrackingStatusShipmentInfoReceived
	case "GWMS_PACKAGE":
		return core.TrackingStatusPackagingComplete
	case "GWMS_OUTBOUND":
		return core.TrackingStatusDispatchedFromWarehouse
	case "WMS_CONFIRMED":
		return core.TrackingStatusWMSConfirmed
	case "SC_INBOUND_SUCCESS":
		return core.TrackingStatusArrivedAtSortingCenter
	case "PU_PICKUP_SUCCESS":
		return core.TrackingStatusAcceptedByCarrier
	case "SC_OUTBOUND_SUCCESS":
		return core.TrackingStatusDepartedFromSortingCenter
	case "LH_HO_IN_SUCCESS":
		return core.TrackingStatusArrivedAtDepartureHub
	case "TRANSIT_PORT_REROUTE_CALLBACK":
		return core.TrackingStatusTransitPortRerouteCb
	case "CC_EX_START":
		return core.TrackingStatusExportCustomsClearanceStarted
	case "CC_EX_SUCCESS":
		return core.TrackingStatusExportCustomsClearanceSuccess
	case "LH_HO_AIRLINE":
		return core.TrackingStatusLeavignDepartureRegion
	case "CC_IM_START":
		return core.TrackingStatusImportCustomsClearanceStarted
	case "LH_DEPART":
		return core.TrackingStatusDepartedOriginRegion
	case "LH_ARRIVE":
		return core.TrackingStatusArrivedAtLinehaulOffice
	case "CC_HO_IN_SUCCESS":
		return core.TrackingStatusArrivedAtCustoms
	case "CC_HO_OUT_SUCCESS":
		return core.TrackingStatusDepartedFromCustoms
	case "CC_IM_SUCCESS":
		return core.TrackingStatusImportCustomsClearanceSuccess
	case "CUSTOMS_ARRIVED_IN_AREA_CALLBACK":
		return core.TrackingStatusArrivedAtCustoms
	default:
		return core.TrackingStatusUnknown
	}
}
