package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

var book_schema = `
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    publishedyear INT NOT NULL
)`

var task_schema = `
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL,
    createdat TIMESTAMP NOT NULL 
)`

func Connect() {
	instance, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	instance.SetMaxOpenConns(10)
	instance.SetMaxIdleConns(10)
	instance.MustExec(book_schema) // Creates books table if not exists
	instance.MustExec(task_schema) // Creates tasks table if not exists
	DB = instance
	fmt.Println("Database connected and table ensured.")
}
