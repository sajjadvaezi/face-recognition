package internal

import (
	"net/http"
)

// SetupRouter initializes the router and defines the routes
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Define routes for face registration and recognition
	mux.HandleFunc("/register", RegisterHandler)   // Handles POST requests to register faces
	mux.HandleFunc("/recognize", RecognizeHandler) // Handles POST requests to recognize faces

	return mux
}
