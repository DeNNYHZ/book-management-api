package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go-fiber-postgres/models"
	"go-fiber-postgres/storage"
	"gorm.io/gorm"
)

// Repository contains database interactions
type Repository struct {
	DB *gorm.DB
}

// CreateBook creates a new book in the database
func (r *Repository) CreateBook(c *fiber.Ctx) error {
	var books []models.Book

	// Parse the request body to handle multiple books
	if err := c.BodyParser(&books); err != nil {
		logrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validate each book
	for i, book := range books {
		if err := book.Validate(); err != nil {
			logrus.WithError(err).Error("Validation failed for book", i)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Insert books into the database
	if err := r.DB.Create(&books).Error; err != nil {
		logrus.WithError(err).Error("Failed to create books")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create books",
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Books created successfully",
		"data":    books,
	})
}

// DeleteBook deletes a book by ID from the database
func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := r.DB.Delete(&models.Book{}, id).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}

// GetBooks retrieves all books from the database with pagination
func (r *Repository) GetBooks(c *fiber.Ctx) error {
	// Retrieve query parameters
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	// Convert page and limit from strings to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logrus.WithError(err).Error("Invalid page number")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logrus.WithError(err).Error("Invalid limit number")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit number",
		})
	}

	// Ensure page and limit are positive
	if page < 1 || limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Page and limit must be greater than 0",
		})
	}

	var books []models.Book
	var totalCount int64

	// Get the total count of books
	if err := r.DB.Model(&models.Book{}).Count(&totalCount).Error; err != nil {
		logrus.WithError(err).Error("Failed to get total count of books")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books count",
		})
	}

	// Fetch books from the database with pagination
	if err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&books).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch books")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Books fetched successfully",
		"data":    books,
		"total":   totalCount,
	})
}

// GetBookByID retrieves a single book by ID from the database
func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var book models.Book

	if err := r.DB.First(&book, id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		logrus.WithError(err).Error("Failed to fetch book")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch book",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book fetched successfully",
		"data":    book,
	})
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	logrus.WithError(err).Error("Internal server error")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
	})
}

// Initialize routes for the API
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// Database configuration
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Initialize database connection
	db, err := storage.NewConnection(config)
	if err != nil {
		logrus.Fatalf("Failed to connect to the database: %v", err)
	}

	// Migrate database schema
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repository and app
	repo := Repository{DB: db}
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust as necessary
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Setup routes
	repo.SetupRoutes(app)

	// Start the server
	if err := app.Listen(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
