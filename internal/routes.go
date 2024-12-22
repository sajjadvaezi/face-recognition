package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"net/http"
	"time"
)

// SetupRouter initializes the router and defines the routes
func SetupRouter() *fiber.App {

	app := fiber.New()

	// Middleware for logging, recovery, and security
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New()) // Use helmet for security headers

	// Middleware for CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust for production to specific origins
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Health check route
	mux := http.NewServeMux()
	mux.HandleFunc("/checkhealth", CheckHealthHandler)
	app.Use(adaptor.HTTPHandler(mux))

	// Serve static files from the "static" directory
	app.Static("/", "./static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 5 * time.Second,
		MaxAge:        3600,
	})

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
	app.Get("/recognize", adaptor.HTTPHandlerFunc(RecognizeHandler))
	app.Post("/add/face", adaptor.HTTPHandlerFunc(AddFaceHandler))
	app.Post("/upload", adaptor.HTTPHandlerFunc(RecognizeWithImageHandler))
	app.Post("/face", adaptor.HTTPHandlerFunc(RegisterFaceWithImageHandler))
	app.Post("/add/class", adaptor.HTTPHandlerFunc(AddClassHandler))
	app.Post("/class/attend", adaptor.HTTPHandlerFunc(AttendanceHandler))
	app.Get("/class/:classname", AttendedUsersHandler)

	return app
}
