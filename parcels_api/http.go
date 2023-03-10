package parcels_api

import (
	"encoding/json"
	"net/http"

	"github.com/dir01/parcels/parcels_service"
	"go.uber.org/zap"
)

func NewServer(parcelsService parcels_service.Service, logger *zap.Logger) *HttpServer {
	return &HttpServer{
		parcelsService: parcelsService,
		logger:         logger,
	}
}

type HttpServer struct {
	parcelsService parcels_service.Service
	logger         *zap.Logger
}

func (s *HttpServer) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/trackingInfo/", s.handleGetTrackingInfo)
	return mux
}

func (s *HttpServer) handleGetTrackingInfo(w http.ResponseWriter, r *http.Request) {
	trackingNumber := r.URL.Query().Get("trackingNumber")
	if trackingNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "message":"trackingNumber query param is required"}`))
		return
	}

	trackingInfos, err := s.parcelsService.GetTrackingInfo(r.Context(), trackingNumber)
	if err != nil {
		s.logger.Error("failed to get tracking info", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message":"internal server error"}`))
		return
	}
	if len(trackingInfos) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"error", "message":"tracking info not found"}`))
		return
	}

	var httpTrackingInfos []*TrackingInfo
	for _, t := range trackingInfos {
		httpTrackingInfos = append(httpTrackingInfos, TrackingInfo{}.fromBusinessStruct(t))
	}

	if respBytes, err := json.Marshal(httpTrackingInfos); err != nil {
		s.logger.Error("failed to marshal response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message":"internal server error"}`))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}
