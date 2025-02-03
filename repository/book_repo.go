package repository

import (
	"coding_test/database"
	"coding_test/models"
)

type BookRepo struct {
	Name string
}

func GetBookRepository() BookRepo {
	return BookRepo{Name: "books"}
}

func (repo BookRepo) Create(book models.Book) error {
	query := `INSERT INTO books (id, title, author, publishedyear) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, book.ID, book.Title, book.Author, book.PublishedYear)
	return err
}

func (repo BookRepo) GetAll() ([]models.Book, error) {
	books := []models.Book{}
	err := database.DB.Select(&books, "SELECT * FROM books")
	return books, err
}

func (repo BookRepo) GetById(book models.Book, bookID int) error {
	query := "SELECT title, author, publishedyear FROM books WHERE id = $1"
	err := database.DB.Get(&book, query, bookID)
	return err
}

func (repo BookRepo) Delete(bookID int) (int64, error) {
	query := "DELETE FROM books WHERE id = $1"
	result, err := database.DB.Exec(query, bookID)
	if err != nil {

	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

func (repo BookRepo) Update(book models.Book, bookID int) error {
	query := "UPDATE books SET title = $1, author = $2, publishedyear = $3 WHERE id = $4"
	_, err := database.DB.Exec(query, book.Title, book.Author, book.PublishedYear, bookID)
	return err
}
