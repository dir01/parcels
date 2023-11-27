package cainiao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dir01/parcels/service"
)

const APIName service.APIName = "cainiao"

func New() service.PostalAPI {
	return &Cainiao{}
}

type Cainiao struct{}

func (c *Cainiao) Fetch(ctx context.Context, trackingNumber string) service.PostalApiResponse {
	result := service.PostalApiResponse{
		TrackingNumber: trackingNumber,
		APIName:        "cainiao",
	}

	url := fmt.Sprintf("https://global.cainiao.com/global/detail.json?mailNos=%s&lang=en-US", trackingNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Status = service.StatusUnknownError
		return result
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		result.Status = service.StatusUnknownError
		return result
	}

	if resp.StatusCode == http.StatusNotFound {
		result.Status = service.StatusNotFound
		return result
	}

	if resp.StatusCode != http.StatusOK {
		result.Status = service.StatusUnknownError
		if b, err := io.ReadAll(resp.Body); err == nil {
			result.ResponseBody = b
		}
		return result
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = service.StatusUnknownError
		return result
	} else {
		result.ResponseBody = responseBody
	}

	if bytes.Contains(responseBody, []byte(`"detailList":[]`)) { // cainiao returns 200 with empty detailList when tracking number is not found
		result.Status = service.StatusNotFound
		return result
	}

	result.Status = service.StatusSuccess

	return result
}

func (c *Cainiao) Parse(rawResponse service.PostalApiResponse) (*service.TrackingInfo, error) {
	var cainiaoResponse response
	if err := json.Unmarshal(rawResponse.ResponseBody, &cainiaoResponse); err != nil {
		return nil, err
	}

	if len(cainiaoResponse.Module) != 1 {
		return nil, fmt.Errorf("response module length is not 1")
	}

	m0 := cainiaoResponse.Module[0]
	var events []service.TrackingEvent
	for _, detail := range m0.DetailList {
		if trackingEvent := c.parseDetail(detail); trackingEvent != nil {
			events = append(events, *trackingEvent)
		}
	}

	return &service.TrackingInfo{
		TrackingNumber:     rawResponse.TrackingNumber,
		APIName:            APIName,
		OriginCountry:      m0.OriginCountry,
		DestinationCountry: m0.DestCountry,
		Events:             events,
	}, nil
}

func (c *Cainiao) parseDetail(detail detail) *service.TrackingEvent {
	return &service.TrackingEvent{
		Time:        time.Unix(detail.Time/1000, 0),
		Description: detail.StanderdDesc,
		Status:      c.mapStatus(detail.ActionCode),
	}
}

func (c *Cainiao) mapStatus(actionCode string) service.TrackingStatus {
	switch actionCode {
	case "GWMS_ACCEPT":
		return service.TrackingStatusShipmentInfoReceived
	case "GWMS_PACKAGE":
		return service.TrackingStatusPackagingComplete
	case "GWMS_OUTBOUND":
		return service.TrackingStatusDispatchedFromWarehouse
	case "WMS_CONFIRMED":
		return service.TrackingStatusWMSConfirmed
	case "SC_INBOUND_SUCCESS":
		return service.TrackingStatusArrivedAtSortingCenter
	case "PU_PICKUP_SUCCESS":
		return service.TrackingStatusAcceptedByCarrier
	case "SC_OUTBOUND_SUCCESS":
		return service.TrackingStatusDepartedFromSortingCenter
	case "LH_HO_IN_SUCCESS":
		return service.TrackingStatusArrivedAtDepartureHub
	case "TRANSIT_PORT_REROUTE_CALLBACK":
		return service.TrackingStatusTransitPortRerouteCb
	case "CC_EX_START":
		return service.TrackingStatusExportCustomsClearanceStarted
	case "CC_EX_SUCCESS":
		return service.TrackingStatusExportCustomsClearanceSuccess
	case "LH_HO_AIRLINE":
		return service.TrackingStatusLeavignDepartureRegion
	case "CC_IM_START":
		return service.TrackingStatusImportCustomsClearanceStarted
	case "LH_DEPART":
		return service.TrackingStatusDepartedOriginRegion
	case "LH_ARRIVE":
		return service.TrackingStatusArrivedAtLinehaulOffice
	case "CC_HO_IN_SUCCESS":
		return service.TrackingStatusArrivedAtCustoms
	case "CC_HO_OUT_SUCCESS":
		return service.TrackingStatusDepartedFromCustoms
	case "CC_IM_SUCCESS":
		return service.TrackingStatusImportCustomsClearanceSuccess
	case "CUSTOMS_ARRIVED_IN_AREA_CALLBACK":
		return service.TrackingStatusArrivedAtCustoms
	default:
		return service.TrackingStatusUnknown
	}
}
