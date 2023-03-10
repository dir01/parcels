package sqlite_storage

import (
	"time"

	"github.com/dir01/parcels/parcels_service"
)

type DBRawPostalApiResponse struct {
	ID             int64  `db:"id"`
	ApiName        string `db:"api_name"`
	TrackingNumber string `db:"tracking_number"`
	FirstFetchedAt int64  `db:"first_fetched_at"`
	LastFetchedAt  int64  `db:"last_fetched_at"`
	ResponseBody   []byte `db:"response_body"`
	Status         string `db:"status"`
}

func (r DBRawPostalApiResponse) ToBusinessModel() *parcels_service.PostalApiResponse {
	return &parcels_service.PostalApiResponse{
		ID:             r.ID,
		ApiName:        r.ApiName,
		TrackingNumber: r.TrackingNumber,
		FirstFetchedAt: fromUnixTime(r.FirstFetchedAt),
		LastFetchedAt:  fromUnixTime(r.LastFetchedAt),
		ResponseBody:   r.ResponseBody,
		Status:         parcels_service.ApiResponseStatus(r.Status),
	}
}

func (r DBRawPostalApiResponse) FromBusinessModel(rawResp *parcels_service.PostalApiResponse) *DBRawPostalApiResponse {
	r.ID = rawResp.ID
	r.ApiName = rawResp.ApiName
	r.TrackingNumber = rawResp.TrackingNumber
	r.FirstFetchedAt = toUnixTime(rawResp.FirstFetchedAt)
	r.LastFetchedAt = toUnixTime(rawResp.LastFetchedAt)
	r.ResponseBody = rawResp.ResponseBody
	r.Status = string(rawResp.Status)
	return &r
}

func toUnixTime(t time.Time) int64 {
	return t.Unix()
}

func fromUnixTime(t int64) time.Time {
	return time.Unix(t, 0)
}
