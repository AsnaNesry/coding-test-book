package handlers

import (
	"coding_test/models"
	"coding_test/repository"
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
	bookRepo := repository.GetBookRepository()
	err := bookRepo.Create(book)
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

	bookRepo := repository.GetBookRepository()
	books, err := bookRepo.GetAll()
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
	bookRepo := repository.GetBookRepository()
	err := bookRepo.GetById(book, bookID)
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
	bookRepo := repository.GetBookRepository()
	rowsAffected, err := bookRepo.Delete(bookID)

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
	bookRepo := repository.GetBookRepository()
	err := bookRepo.Update(book, bookID)
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
