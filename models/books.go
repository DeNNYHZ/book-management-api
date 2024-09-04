package models

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Book represents a book record
type Book struct {
	gorm.Model
	Author    string `json:"author" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
}

// Initialize the validator
var validate = validator.New()

// Validate checks the model's fields
func (b *Book) Validate() error {
	return validate.Struct(b)
}

// MigrateBooks migrates the book schema
func MigrateBooks(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}

// MarshalJSON customizes the JSON representation of Book
func (b *Book) MarshalJSON() ([]byte, error) {
	type Alias Book
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
		Alias:     (*Alias)(b),
	})
}
