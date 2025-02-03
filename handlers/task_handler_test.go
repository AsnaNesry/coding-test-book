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

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {

	db := setupTestDB(t)
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

	db := setupTestDB(t)
	defer db.Close()

	database.DB = db
	date := time.Now().UTC()
	formattedDate := date.Format("2006-01-02T15:04:05.999999Z")

	_, err := database.DB.Exec(`INSERT INTO tasks (id, title, completed, createdat) VALUES ($1, $2, $3, $4)`,
		3, "Task2", false, formattedDate)
	if err != nil {
		t.Fatalf("Error inserting test task: %v", err)
	}

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
