package main

import (
	"github.com/dir01/parcels/metrics"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dir01/parcels/externalapis/cainiao"
	"github.com/dir01/parcels/parcels_api"
	"github.com/dir01/parcels/service"
	"github.com/dir01/parcels/sqlite_storage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func main() {
	// region parameters
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		panic("DB_PATH is not set")
	}

	bindAddr := "0.0.0.0:0"
	if bindAddrEnv := os.Getenv("BIND_ADDR"); bindAddrEnv != "" {
		bindAddr = bindAddrEnv
	}

	okCheckInterval := 24 * time.Hour           // how often to check after a successful fetch
	notFoundCheckInterval := 3 * 24 * time.Hour // how often to check after a not found response
	unknownErrorCheckInterval := 3 * time.Hour  // how often to check after an unknown error
	apiFetchTimeout := 10 * time.Second         // how long to wait for a response from an API
	// expiryTimeout is the time after which a parcel is treated as if we never heard of it
	// this is due to the fact that sometimes tracking numbers can be reused
	expiryTimeout := 6 * 30 * 24 * time.Hour
	// endregion

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	db := sqlx.MustOpen("sqlite3", dbPath)
	storage := sqlite_storage.NewStorage(db)

	promMetrics := metrics.NewPrometheus([]service.APIName{cainiao.APIName})

	apiMap := map[service.APIName]service.PostalAPI{
		cainiao.APIName: cainiao.New(),
	}

	svc := service.NewService(
		apiMap,
		storage,
		promMetrics,
		okCheckInterval,
		notFoundCheckInterval,
		unknownErrorCheckInterval,
		apiFetchTimeout,
		expiryTimeout,
		logger,
		time.Now,
	)

	httpServer := parcels_api.NewServer(svc, logger)
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic(err)
	}
	logger.Info("listening", zap.String("addr", listener.Addr().String()))
	err = http.Serve(listener, httpServer.GetMux())
	logger.Info("svc terminated", zap.Error(err))
}
