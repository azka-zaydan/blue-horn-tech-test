package main

import (
	"log"
	"os"
	"time"

	"mini-evv-logger-backend/config"
	"mini-evv-logger-backend/src/domains/schedule/controller"
	scheduleRepo "mini-evv-logger-backend/src/domains/schedule/repository"
	scheduleService "mini-evv-logger-backend/src/domains/schedule/service"
	taskController "mini-evv-logger-backend/src/domains/task/controller"
	taskRepo "mini-evv-logger-backend/src/domains/task/repository"
	taskService "mini-evv-logger-backend/src/domains/task/service"
	"mini-evv-logger-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Zerolog
	utils.InitLogger()
	mainLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	mainLogger.Info().Msg("Starting Mini EVV Logger Backend")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	db, err := config.InitDB(cfg, mainLogger)
	if err != nil {
		mainLogger.Fatal().Err(err).Msg("Failed to initialize database connection")
	}
	defer db.Close()

	mainLogger.Info().Msg("Successfully connected to PostgreSQL database")

	// Initialize Repositories (now returning interfaces)
	scheduleRepository := scheduleRepo.NewScheduleRepository(db, mainLogger)
	taskRepository := taskRepo.NewTaskRepository(db, mainLogger)

	// Initialize Services (now returning interfaces)
	// Now injecting taskRepository directly into NewScheduleService
	scheduleSvc := scheduleService.NewScheduleService(scheduleRepository, taskRepository)
	taskSvc := taskService.NewTaskService(taskRepository)

	// Initialize Controllers (now injecting service interfaces)
	scheduleCtrl := controller.NewScheduleController(scheduleSvc)
	taskCtrl := taskController.NewTaskController(taskSvc)

	// Initialize Fiber app
	app := fiber.New()

	// Add Fiber Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} ${locals:requestid}\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	// Apply CORS middleware to allow cross-origin requests
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           // Allows all origins, you can restrict this to specific origins (e.g., "http://localhost:3000")
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",                 // Allowed HTTP methods
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Allowed headers
	}))

	// Basic root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Welcome to Mini EVV Logger API!",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Define API group
	api := app.Group("/api")

	// Register routes using controller methods
	scheduleCtrl.Routes(api)
	taskCtrl.Routes(api)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	mainLogger.Info().Msgf("Server listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		mainLogger.Fatal().Err(err).Msg("Failed to start server")
	}
}
