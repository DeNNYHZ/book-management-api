package main

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go-fiber-postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Helper function to set up an in-memory SQLite database
func setupTestDB() *gorm.DB {
	dsn := "user=postgres password=admin dbname=fiber_demo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Book{})
	return db
}

func TestCreateBook(t *testing.T) {
	app := fiber.New()
	repo := &Repository{DB: setupTestDB()}
	repo.SetupRoutes(app)

	reqBody := `[{"author":"Test Author","title":"Test Title","publisher":"Test Publisher"}]`
	req := httptest.NewRequest("POST", "/api/create_books", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestDeleteBook(t *testing.T) {
	app := fiber.New()
	repo := &Repository{DB: setupTestDB()}
	app.Delete("/api/delete_book/:id", repo.DeleteBook)

	// Create a book to delete
	book := models.Book{Author: "Test Author", Title: "Test Title", Publisher: "Test Publisher"}
	repo.DB.Create(&book)

	// Test deleting the book
	req := httptest.NewRequest("DELETE", "/api/delete_book/"+strconv.Itoa(int(book.ID)), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check if the book was actually deleted
	var count int64
	repo.DB.Model(&models.Book{}).Where("id = ?", book.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestGetBooks(t *testing.T) {
	app := fiber.New()
	repo := &Repository{DB: setupTestDB()}
	app.Get("/api/books", repo.GetBooks)

	// Add some books to the database
	for i := 1; i <= 5; i++ {
		repo.DB.Create(&models.Book{
			Author:    "Author " + strconv.Itoa(i),
			Title:     "Title " + strconv.Itoa(i),
			Publisher: "Publisher " + strconv.Itoa(i),
		})
	}

	req := httptest.NewRequest("GET", "/api/books?page=1&limit=10", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGetBookByID(t *testing.T) {
	app := fiber.New()
	repo := &Repository{DB: setupTestDB()}
	app.Get("/api/get_books/:id", repo.GetBookByID)

	// Create a book to retrieve
	book := models.Book{Author: "Test Author", Title: "Test Title", Publisher: "Test Publisher"}
	repo.DB.Create(&book)

	req := httptest.NewRequest("GET", "/api/get_books/"+strconv.Itoa(int(book.ID)), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
