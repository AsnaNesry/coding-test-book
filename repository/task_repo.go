package repository

import (
	"coding_test/database"
	"coding_test/models"
)

type TaskRepo struct {
	Name string
}

func GetTaskRepository() TaskRepo {
	return TaskRepo{Name: "tasks"}
}

func (taskRepo TaskRepo) Create(task models.Task) error {
	query := `INSERT INTO tasks (id, title, completed) VALUES ($1, $2, $3)`
	_, err := database.DB.Exec(query, task.ID, task.Title, task.Completed)
	return err
}

func (taskRepo TaskRepo) GetAll(tasks []models.Task) error {
	err := database.DB.Select(&tasks, "SELECT * FROM tasks")
	return err
}

func (taskRepo TaskRepo) Delete(taskID int) (int64, error) {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := database.DB.Exec(query, taskID)
	if err != nil {
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected, err
}

func (taskRepo TaskRepo) MarkTaskCompleted(id int) (int64, error) {
	query := "UPDATE tasks SET completed = true WHERE id = $1"
	queryResult, err := database.DB.Exec(query, id)
	updateCount, _ := queryResult.RowsAffected()
	return updateCount, err
}

func (taskRepo TaskRepo) Update(task models.Task, taskID int) error {
	query := "UPDATE tasks SET title = $1, completed = $2 WHERE id = $3"
	_, err := database.DB.Exec(query, task.Title, task.Completed, taskID)
	return err
}
