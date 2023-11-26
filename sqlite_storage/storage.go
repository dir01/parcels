package sqlite_storage

import (
	"context"
	"strings"

	"github.com/dir01/parcels/service"
	"github.com/hori-ryota/zaperr"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func NewStorage(db *sqlx.DB) service.Storage {
	return &sqliteStorage{db: db}
}

type sqliteStorage struct {
	db *sqlx.DB
}

func (s sqliteStorage) GetLatest(
	ctx context.Context,
	trackingNumber string,
	apiNames []service.APIName,
) ([]*service.PostalApiResponse, error) {
	var apiNamesStr []string
	for _, apiName := range apiNames {
		apiNamesStr = append(apiNamesStr, string(apiName))
	}
	zapFields := []zap.Field{
		zap.String("trackingNumber", trackingNumber),
		zap.Strings("apiNames", apiNamesStr),
	}
	rows, err := s.db.QueryxContext(ctx, `
		SELECT * FROM postal_api_responses 
		         WHERE tracking_number = ? 
		           AND api_name IN (?) 
		         ORDER BY last_fetched_at DESC;`, trackingNumber, strings.Join(apiNamesStr, ","))
	if err != nil {
		return nil, zaperr.Wrap(err, "failed to QueryxContext", zapFields...)
	}
	defer rows.Close()

	var dbStructs []DBRawPostalApiResponse
	for rows.Next() {
		var dbStruct DBRawPostalApiResponse
		err := rows.StructScan(&dbStruct)
		if err != nil {
			return nil, zaperr.Wrap(err, "failed to StructScan", zapFields...)
		}
		dbStructs = append(dbStructs, dbStruct)
	}

	var businessStructs []*service.PostalApiResponse
	for _, dbStruct := range dbStructs {
		businessStructs = append(businessStructs, dbStruct.ToBusinessModel())
	}

	return businessStructs, nil
}

func (s sqliteStorage) Insert(ctx context.Context, trackingNumber string, apiName service.APIName, response *service.PostalApiResponse) error {
	dbStruct := DBRawPostalApiResponse{}.FromBusinessModel(response)
	dbStruct.TrackingNumber = trackingNumber
	dbStruct.APIName = apiName
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO postal_api_responses 
		    (api_name, tracking_number, first_fetched_at, last_fetched_at, response_body, status)
		VALUES 
		    (:api_name, :tracking_number, :first_fetched_at, :last_fetched_at, :response_body, :status)
	`, dbStruct)

	if err != nil {
		return zaperr.Wrap(err, "failed to NamedExecContext", zap.Any("dbStruct", dbStruct))
	}
	return nil
}

func (s sqliteStorage) Update(ctx context.Context, response *service.PostalApiResponse) error {
	dbStruct := DBRawPostalApiResponse{}.FromBusinessModel(response)
	_, err := s.db.NamedExecContext(ctx, `
		UPDATE postal_api_responses
		SET api_name = :api_name,
		    tracking_number = :tracking_number,
		    first_fetched_at = :first_fetched_at,
		    last_fetched_at = :last_fetched_at,
		    response_body = :response_body,
		    status = :status
		WHERE id = :id
	`, dbStruct)

	if err != nil {
		return zaperr.Wrap(err, "failed to NamedExecContext", zap.Any("dbStruct", dbStruct))
	}
	return nil
}
