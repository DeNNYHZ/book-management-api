package Testing

import (
	"bytes"
	"gorm.io/driver/sqlite"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go-fiber-postgres/models"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&models.Book{})
	return db
}

type Repository struct {
	DB *gorm.DB
}

// Initialize routes for the API
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books")
	api.Delete("/delete_book/:id")
	api.Get("/get_books/:id")
	api.Get("/books")
}

func TestCreateBook(t *testing.T) {
	app := fiber.New()
	repo := &Repository{DB: setupTestDB()}
	repo.SetupRoutes(app)

	reqBody := `{"author":"Test Author","title":"Test Title","publisher":"Test Publisher"}`
	req := httptest.NewRequest("POST", "/api/create_books", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
