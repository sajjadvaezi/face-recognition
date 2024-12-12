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

	port := ":8090"

	err := fiberApp.Listen(port)
	if err != nil {
		logger.Error(fmt.Sprintf("could not start app on port %s error: %s", port, err.Error()))
	}

}
