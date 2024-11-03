package main

import (
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal"

	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db.InitSQLite()

	mux := internal.SetupRouter()

	err := http.ListenAndServe("localhost:8090", mux)
	if err != nil {
		return
	}
}
