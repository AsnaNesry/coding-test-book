package main

import (
	"coding_test/database"
	handlers "coding_test/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	r.HandleFunc("/books", handlers.GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.GetBookByID).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.UpdateBookByID).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteBookByID).Methods("DELETE")

	r.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", handlers.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.UpdateTaskByID).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTaskByID).Methods("DELETE")
	r.HandleFunc("/tasks/bulk", handlers.UpdateAllTask).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", r))
}
