package sqlite_storage

import (
	"time"

	"github.com/dir01/parcels/service"
)

type DBRawPostalApiResponse struct {
	ID             int64           `db:"id"`
	APIName        service.APIName `db:"api_name"`
	TrackingNumber string          `db:"tracking_number"`
	FirstFetchedAt int64           `db:"first_fetched_at"`
	LastFetchedAt  int64           `db:"last_fetched_at"`
	ResponseBody   []byte          `db:"response_body"`
	Status         string          `db:"status"`
}

func (r DBRawPostalApiResponse) ToBusinessModel() *service.PostalApiResponse {
	return &service.PostalApiResponse{
		ID:             r.ID,
		APIName:        r.APIName,
		TrackingNumber: r.TrackingNumber,
		FirstFetchedAt: fromUnixTime(r.FirstFetchedAt),
		LastFetchedAt:  fromUnixTime(r.LastFetchedAt),
		ResponseBody:   r.ResponseBody,
		Status:         service.ApiResponseStatus(r.Status),
	}
}

func (r DBRawPostalApiResponse) FromBusinessModel(rawResp *service.PostalApiResponse) *DBRawPostalApiResponse {
	r.ID = rawResp.ID
	r.APIName = rawResp.APIName
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
