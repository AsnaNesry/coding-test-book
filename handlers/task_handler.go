package handlers

import (
	"coding_test/database"
	"coding_test/models"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// funtion to create a task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, task.ID, task.Title, task.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "Task created successfully",
		"id":      task.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// function to read all the tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	err := database.DB.Select(&tasks, "SELECT * FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// function to delete a task
func DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	taskID, _ := strconv.Atoi(id)
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := database.DB.Exec(query, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted"})
}

// function to update a task details
func UpdateTaskByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	taskID, _ := strconv.Atoi(id)

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE tasks SET title = $1, completed = $2 WHERE id = $3"
	_, err := database.DB.Exec(query, task.Title, task.Completed, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Task updated successfully",
		"id":      taskID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// function to update task in bulk

type QueryResult struct {
	Id           int
	RowsAffected int64
	Err          error
}

func UpdateAllTask(w http.ResponseWriter, r *http.Request) {

	var ids []int
	var wg sync.WaitGroup
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	ch := make(chan QueryResult)
	wg.Add(len(ids))
	for _, id := range ids {
		go markTaskCompleted(&wg, id, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	isSuccess := true
	var failedIds []int
	for res := range ch {
		if res.Err != nil {
			isSuccess = false
			failedIds = append(failedIds, res.Id)
		}
	}
	var response map[string]interface{}
	if isSuccess {
		response = map[string]interface{}{
			"message": "All tasks marked completed",
		}
		w.WriteHeader(http.StatusOK)
	} else {
		response = map[string]interface{}{
			"message": "Failed to mark completed ",
			"Id":      failedIds,
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func markTaskCompleted(wg *sync.WaitGroup, id int, ch chan<- QueryResult) {
	defer wg.Done()
	query := "UPDATE tasks SET completed = true WHERE id = $1"
	queryResult, err := database.DB.Exec(query, id)
	updateCount, _ := queryResult.RowsAffected()
	ch <- QueryResult{Id: id, RowsAffected: updateCount, Err: err}
}
