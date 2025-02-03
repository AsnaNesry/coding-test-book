package handlers

import (
	"bytes"
	"coding_test/database"
	"coding_test/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")

	if err != nil {
		t.Fatalf("Error setting up the test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
    				id SERIAL PRIMARY KEY,
    				title TEXT NOT NULL,
    				author TEXT NOT NULL,
    				publishedyear INT NOT NULL
					)`)
	if err != nil {
		t.Fatalf("Error creating test table: %v", err)
	}
	db.Exec("DELETE FROM books")

	return db
}

func TestCreateBook(t *testing.T) {

	db := setupTestDB(t)
	defer db.Close()

	database.DB = db
	handler := http.HandlerFunc(CreateBook)

	bookJson := getBookJson(t)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(bookJson))

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status 201 Created")
}

func TestGetAllBooks(t *testing.T) {

	db := setupTestDB(t)
	defer db.Close()

	database.DB = db

	_, err := database.DB.Exec(`INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)`,
		3, "Book2", "Author2", 2020)
	if err != nil {
		t.Fatalf("Error inserting test book: %v", err)
	}

	handler := http.HandlerFunc(GetAllBooks)
	request := httptest.NewRequest(http.MethodGet, "/books", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
}

func getBookJson(t *testing.T) []byte {
	book := models.Book{
		ID:            1,
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedYear: 2025,
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Error marshalling book: %v", err)
	}
	return bookJson
}
