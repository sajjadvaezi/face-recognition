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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db.InitSQLite()

	fiberApp := internal.SetupRouter()

	// Set port and host to listen on all interfaces
	port := ":8090"

	// Use 0.0.0.0 to listen on all interfaces (including external devices)
	err := fiberApp.Listen("0.0.0.0" + port)
	if err != nil {
		logger.Error(fmt.Sprintf("could not start app on port %s error: %s", port, err.Error()))
	}
}
