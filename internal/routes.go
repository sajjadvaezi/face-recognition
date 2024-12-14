package internal

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"net/http"
)

// SetupRouter initializes the router and defines the routes
func SetupRouter() *fiber.App {

	app := fiber.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/checkhealth", CheckHealthHandler)

	// Define routes for face registration and recognition

	app.Post("/register", adaptor.HTTPHandlerFunc(RegisterHandler))
	// Handles POST requests to recognize faces

	app.Get("/recognize", adaptor.HTTPHandlerFunc(RecognizeHandler))
	app.Post("/add/face", adaptor.HTTPHandlerFunc(AddFaceHandler))
	app.Post("/upload", adaptor.HTTPHandlerFunc(RecognizeWithImageHandler))

	// this is for adding face with request containing image in it

	app.Post("/face", adaptor.HTTPHandlerFunc(RegisterFaceWithImageHandler))
	app.Post("/add/class", adaptor.HTTPHandlerFunc(AddClassHandler))

	app.Post("/class/attend", adaptor.HTTPHandlerFunc(AttendanceHandler))

	app.Get("/class/:classname", AttendedUsersHandler)

	return app
}
