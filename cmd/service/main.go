package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dir01/parcels/core"
	"github.com/dir01/parcels/http_server"
	"github.com/dir01/parcels/postal_apis/cainiao"
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

	okCheckInterval := 24 * time.Hour
	notFoundCheckInterval := 3 * 24 * time.Hour
	unknownErrorCheckInterval := 3 * time.Hour
	apiFetchTimeout := 10 * time.Second
	expiryTimeout := 6 * 30 * 24 * time.Hour
	// endregion

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db := sqlx.MustOpen("sqlite3", dbPath)
	storage := sqlite_storage.NewStorage(db)

	apiMap := map[string]core.PostalApi{
		"cainiao": cainiao.New(),
	}

	service := core.NewService(
		apiMap,
		storage,
		okCheckInterval,
		notFoundCheckInterval,
		unknownErrorCheckInterval,
		apiFetchTimeout,
		expiryTimeout,
		logger,
		time.Now,
	)

	httpServer := http_server.NewServer(service, logger)
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic(err)
	}
	logger.Info("listening", zap.String("addr", listener.Addr().String()))
	err = http.Serve(listener, httpServer.GetMux())
	logger.Info("service terminated", zap.Error(err))
}
