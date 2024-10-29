package main

import (
	"database/sql"
	"github.com/sajjadvaezi/face-recognition/internal"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a database connection
	db, err := sql.Open("sqlite3", "./face_recognition.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := internal.SetupRouter()
	
	err = http.ListenAndServe("localhost:8090", mux)
	if err != nil {
		return
	}
}
