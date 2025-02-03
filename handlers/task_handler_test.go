package handlers

import (
	"bytes"
	"coding_test/database"
	"coding_test/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDBForTasks(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")

	if err != nil {
		t.Fatalf("Error setting up the test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
						id SERIAL PRIMARY KEY,
						title TEXT NOT NULL,
						completed BOOLEAN NOT NULL,
						createdat TIMESTAMP NOT NULL DEFAULT NOW()
					)`)
	if err != nil {
		t.Fatalf("Error creating test table: %v", err)
	}
	db.Exec("DELETE FROM tasks")

	return db
}

func TestBulkUpdate(t *testing.T) {
	db := setupTestDBForTasks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "1", "Title1", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "2", "Title2", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "3", "Title3", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "4", "Title4", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "5", "Title5", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "6", "Title6", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "7", "Title7", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "8", "Title8", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "9", "Title9", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "10", "Title10", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "11", "Title11", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "12", "Title12", "false")
	tx.Commit()

	handler := http.HandlerFunc(UpdateAllTask)
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	jsonIds, err := json.Marshal(ids)
	if err != nil {
		t.Fatalf("Error marshalling task: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/tasks/bulk", bytes.NewReader(jsonIds))

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status 200 Ok")

	completedCount := 0

	db.Get(&completedCount, "SELECT count(*) FROM tasks where completed = true")

	assert.Equal(t, 12, completedCount, "All tasks should be marked as completed")
}

func TestCreateTask(t *testing.T) {

	db := setupTestDBForTasks(t)
	defer db.Close()

	database.DB = db
	handler := http.HandlerFunc(CreateTask)

	taskJson := getTaskJson(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(taskJson))

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected HTTP status 201 Created")
}

func TestGetAllTasks(t *testing.T) {

	db := setupTestDBForTasks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "1", "Title1", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "2", "Title2", "true")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "3", "Title3", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "12", "Title12", "true")
	tx.Commit()

	handler := http.HandlerFunc(GetAllTasks)
	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")
}

func getTaskJson(t *testing.T) []byte {
	dateString := "2025-02-02T14:37:35.591258Z"

	parsedTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		t.Fatalf("Error parsing time: %v", err)
	}
	task := models.Task{
		ID:        4,
		Title:     "Task6",
		Completed: false,
		CreatedAt: parsedTime,
	}

	bookJson, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Error marshalling task: %v", err)
	}
	return bookJson
}

func TestUpdateTaskById(t *testing.T) {

	db := setupTestDBForTasks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "1", "Title1", "false")
	tx.Commit()
	handler := http.HandlerFunc(UpdateTaskByID)
	request := httptest.NewRequest(http.MethodPut, "/tasks/1", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	completed := false

	db.Get(&completed, "SELECT completed FROM tasks where id = 1")

	assert.Equal(t, true, completed, "Updated values should match")
}

func TestDeleteTaskById(t *testing.T) {

	db := setupTestDBForTasks(t)
	defer db.Close()

	database.DB = db

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "1", "Title1", "false")
	tx.MustExec("INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)", "2", "Title2", "false")
	tx.Commit()

	handler := http.HandlerFunc(DeleteTaskByID)
	request := httptest.NewRequest(http.MethodPut, "/tasks/2", nil)
	request = mux.SetURLVars(request, map[string]string{"id": "2"})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, "Expected HTTP status 200 OK")

	recordCount := 0

	db.Get(&recordCount, "SELECT count(*) FROM tasks")

	assert.Equal(t, 1, recordCount, "Only one record should exist")
}
