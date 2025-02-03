package handlers

import (
	"bytes"
	"coding_test/database"
	"coding_test/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDBForBooks(t *testing.T) *sqlx.DB {
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

	db := setupTestDBForBooks(t)
	defer db.Close()

	database.DB = db
	handler := http.HandlerFunc(CreateBook)

	bookJson := getBookJson(t)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(bookJson))

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status 201 Created")

	bookId := 0

	db.Get(&bookId, "SELECT id FROM books")
	assert.Equal(t, 1, bookId, "Expected Id should be 1")
}

func TestGetAllBooks(t *testing.T) {

	db := setupTestDBForBooks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "1", "Title1", "Author1", "2021")
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "2", "Title2", "Author2", "2022")
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "3", "Title3", "Author3", "2023")
	tx.Commit()

	handler := http.HandlerFunc(GetAllBooks)
	request := httptest.NewRequest(http.MethodGet, "/books", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	existingRecords := 0

	db.Get(&existingRecords, "SELECT count(*) FROM books")

	assert.Equal(t, 3, existingRecords, "Existing records count must match")
}

func TestGetBookById(t *testing.T) {

	db := setupTestDBForBooks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "1", "Title1", "Author1", "2021")
	tx.Commit()

	handler := http.HandlerFunc(GetBookByID)
	request := httptest.NewRequest(http.MethodPut, "/books/1", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	var title string

	db.Get(&title, "SELECT title FROM books where id = 1")

	assert.Equal(t, "Title1", title, "Get By ID, title should match")
}

func TestUpdateBookById(t *testing.T) {

	db := setupTestDBForBooks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "1", "Title1", "Author1", "2021")
	tx.Commit()
	book := models.Book{
		Title:         "Title2",
		Author:        "Author2",
		PublishedYear: 2022,
	}
	bookJson, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Error marshalling book: %v", err)
	}
	handler := http.HandlerFunc(UpdateBookByID)
	request := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(bookJson))
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	var title string

	db.Get(&title, "SELECT title FROM books where id = 1")

	assert.Equal(t, "Title2", title, "Updated values should match")
}

func TestDeleteBookById(t *testing.T) {

	db := setupTestDBForBooks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "1", "Title1", "Author1", "2021")
	tx.MustExec("INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)", "2", "Title2", "Author2", "2022")
	tx.Commit()

	handler := http.HandlerFunc(DeleteBookByID)
	request := httptest.NewRequest(http.MethodPut, "/books/1", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	recordCount := 0

	db.Get(&recordCount, "SELECT count(*) FROM books")

	assert.Equal(t, 1, recordCount, "Only one record should exist")
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
