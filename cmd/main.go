package main

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize SQLite database
	db.InitSQLite()

	// Setup Fiber application
	fiberApp := internal.SetupRouter()

	// Define TLS certificate and key file paths
	certFile := "cert.pem"
	keyFile := "key.pem"

	// Start the server with TLS
	port := ":443"
	logger.Info(fmt.Sprintf("Starting app with TLS on port %s...", port))

	err := fiberApp.ListenTLS(port, certFile, keyFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not start app on port %s. Error: %s", port, err.Error()))
	}
}
