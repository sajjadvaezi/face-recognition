package main

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db.InitSQLite()

	mux := internal.SetupRouter()

	port := "8090"
	addr := fmt.Sprintf("localhost:%s", port)

	logger.Info(fmt.Sprintf("staring server on %s", addr))
	err := http.ListenAndServe(addr, mux)

	if err != nil {
		logger.Error("could not start mux, error:", err)
		return
	}

}
