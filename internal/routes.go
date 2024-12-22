package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/http"
	"time"
)

// SetupRouter initializes the router and defines the routes
func SetupRouter() *fiber.App {

	app := fiber.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/checkhealth", CheckHealthHandler)

	// Serve static files from the "static" directory
	app.Static("/", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 5 * time.Second,
		MaxAge:        3600,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins (adjust this in production for security)
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	// Serve views for registration and face data capture
	app.Get("/register-view", func(c *fiber.Ctx) error {
		return c.SendFile("./views/register.html")
	})

	app.Get("/face-view", func(c *fiber.Ctx) error {
		return c.SendFile("./views/face.html")
	})

	app.Get("/add-class-view", func(c *fiber.Ctx) error {
		return c.SendFile("./views/add_class.html")
	})

	app.Get("/attend-class-view", func(c *fiber.Ctx) error {
		return c.SendFile("./views/attend_class.html")
	})

	app.Get("/show-attendance-view", func(c *fiber.Ctx) error {
		return c.SendFile("./views/show_attendance.html")
	})

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
