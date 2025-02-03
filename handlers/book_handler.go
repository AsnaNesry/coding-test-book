package handlers

import (
	"coding_test/database"
	"coding_test/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// funtion to create a book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, book.ID, book.Title, book.Author, book.PublishedYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Book created successfully",
		"id":      book.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// function to read all the books
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books := []models.Book{}
	err := database.DB.Select(&books, "SELECT * FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// function to get a specific book details w.r.t ID
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	bookID, _ := strconv.Atoi(id)

	var book models.Book
	query := "SELECT title, author, publishedyear FROM books WHERE id = $1"
	err := database.DB.Get(&book, query, bookID)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// function to delete a book
func DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	bookID, _ := strconv.Atoi(id)
	query := "DELETE FROM books WHERE id = $1"
	result, err := database.DB.Exec(query, bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted"})
}

// function to update a book details
func UpdateBookByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	bookID, _ := strconv.Atoi(id)

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE books SET title = $1, author = $2, publishedyear = $3 WHERE id = $4"
	_, err := database.DB.Exec(query, book.Title, book.Author, book.PublishedYear, bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Book updated successfully",
		"id":      bookID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
