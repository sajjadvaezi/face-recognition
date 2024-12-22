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
	// Initialize the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize the SQLite database
	db.InitSQLite()

	// Setup Fiber application
	fiberApp := internal.SetupRouter()

	// Define the certificate and key file paths
	certFile := "cert.pem"
	keyFile := "key.pem"

	// Port for the application
	port := ":8090"

	// Ensure certificate files exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Certificate file %s not found", certFile))
		os.Exit(1)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Key file %s not found", keyFile))
		os.Exit(1)
	}

	// Listen on HTTPS using the provided certificate and key
	err := fiberApp.ListenTLS("0.0.0.0"+port, certFile, keyFile)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not start app on port %s with HTTPS, error: %s", port, err.Error()))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("App is running securely on https://0.0.0.0%s", port))
}
