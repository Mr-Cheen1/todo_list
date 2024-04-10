package db

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID     int       `json:"id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
}

func GetAllTasks(statusFilter, sortOrder string) ([]Task, error) {
	var tasks []Task
	var rows *sql.Rows
	var err error
	query := "SELECT id, task_text, task_date, status FROM tasks"

	if statusFilter != "" {
		query += " WHERE status = '" + statusFilter + "'"
	}

	if sortOrder != "" {
		query += " ORDER BY task_date " + sortOrder
	}

	rows, err = DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if scanErr := rows.Scan(&task.ID, &task.Text, &task.Date, &task.Status); scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func CreateTask(task Task) error {
	query := "INSERT INTO tasks (task_text, task_date, status) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, task.Text, task.Date.UTC(), task.Status)
	return err
}

func UpdateTask(task Task) error {
	result, err := DB.Exec("UPDATE tasks SET task_text = $1, task_date = $2, status = $3 WHERE id = $4",
		task.Text, task.Date, task.Status, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}
