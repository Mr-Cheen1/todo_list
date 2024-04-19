package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	StatusInProgress = iota
	StatusCompleted
)

type Task struct {
	ID           int64     `json:"id"`
	Text         string    `json:"text"`
	CreatedDate  time.Time `json:"createdDate"`
	ExpectedDate time.Time `json:"expectedDate"`
	Status       int       `json:"status"`
}

func (t *Task) MarshalJSON() ([]byte, error) {
	type Alias Task
	return json.Marshal(&struct {
		CreatedDate  string `json:"createdDate"`
		ExpectedDate string `json:"expectedDate"`
		*Alias
	}{
		CreatedDate:  t.CreatedDate.Format("2006-01-02"),
		ExpectedDate: t.ExpectedDate.Format("2006-01-02"),
		Alias:        (*Alias)(t),
	})
}

func GetAllTasks(statusFilter, sortOrder, sortField string) ([]Task, error) {
	var tasks []Task
	var rows *sql.Rows
	var err error

	query := "SELECT id, task_text, createdDate, expectedDate, status FROM tasks"
	var conditions []string
	var args []interface{}

	if statusFilter != "" {
		conditions = append(conditions, "status = $1")
		args = append(args, statusFilter)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Белый список допустимых значений для sortField.
	validSortFields := map[string]bool{
		"id":           true,
		"task_text":    true,
		"createdDate":  true,
		"expectedDate": true,
		"status":       true,
	}

	if sortField != "" {
		// Проверяем, находится ли sortField в белом списке.
		if validSortFields[sortField] {
			query += " ORDER BY " + sortField
			if sortOrder == "desc" {
				query += " DESC"
			}
		} else {
			// Если sortField не находится в белом списке, возвращаем ошибку.
			return nil, fmt.Errorf("invalid sort field: %s", sortField)
		}
	}

	rows, err = DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if scanErr := rows.Scan(&task.ID, &task.Text, &task.CreatedDate, &task.ExpectedDate, &task.Status); scanErr != nil {
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
	query := "INSERT INTO tasks (task_text, createdDate, expectedDate, status) VALUES ($1, $2, $3, $4)"
	_, err := DB.Exec(query, task.Text, task.CreatedDate, task.ExpectedDate, task.Status)
	return err
}

func UpdateTask(task Task) error {
	result, err := DB.Exec(
		"UPDATE tasks SET task_text = $1, createdDate = $2, expectedDate = $3, status = $4 WHERE id = $5",
		task.Text, task.CreatedDate, task.ExpectedDate, task.Status, task.ID,
	)
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

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		CreatedDate  string `json:"createdDate"`
		ExpectedDate string `json:"expectedDate"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	// Обновленный формат для полной даты и времени.
	t.CreatedDate, err = time.Parse(time.RFC3339, aux.CreatedDate)
	if err != nil {
		return err
	}
	t.ExpectedDate, err = time.Parse(time.RFC3339, aux.ExpectedDate)
	if err != nil {
		return err
	}
	return nil
}
