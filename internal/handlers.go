package internal

import (
	"fmt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "register handler")
	if err != nil {
		return
	}
}

func RecognizeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "register handler")
	if err != nil {
		return
	}
}
