package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Константы для статусов задач.
const (
	StatusInProgress = iota
	StatusCompleted
	StatusTesting
	StatusReturned
)

// Структура Task представляет задачу.
type Task struct {
	ID           int64     `json:"id"`
	Text         string    `json:"text"`
	CreatedDate  time.Time `json:"createdDate"`
	ExpectedDate time.Time `json:"expectedDate"`
	Status       int       `json:"status"`
}

// Метод MarshalJSON для сериализации структуры Task в JSON.
func (t *Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":           t.ID,
		"text":         t.Text,
		"createdDate":  t.CreatedDate.Format("2006-01-02"),
		"expectedDate": t.ExpectedDate.Format("2006-01-02"),
		"status":       t.Status,
	})
}

// Функция GetAllTasks получает все задачи из базы данных с учетом фильтрации и сортировки.
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

// Функция CreateTask создает новую задачу в базе данных.
func CreateTask(task Task) error {
	query := "INSERT INTO tasks (task_text, createdDate, expectedDate, status) VALUES ($1, $2, $3, $4)"

	createdDateStr := task.CreatedDate.Format("2006-01-02")
	expectedDateStr := task.ExpectedDate.Format("2006-01-02")
	_, err := DB.Exec(query, task.Text, createdDateStr, expectedDateStr, task.Status)
	return err
}

// Функция UpdateTask обновляет существующую задачу в базе данных.
func UpdateTask(task Task) error {
	createdDateStr := task.CreatedDate.Format("2006-01-02")
	expectedDateStr := task.ExpectedDate.Format("2006-01-02")

	result, err := DB.Exec(
		"UPDATE tasks SET task_text = $1, createdDate = $2, expectedDate = $3, status = $4 WHERE id = $5",
		task.Text, createdDateStr, expectedDateStr, task.Status, task.ID,
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

// Функция DeleteTask удаляет задачу из базы данных по ее идентификатору.
func DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}

// Метод UnmarshalJSON для десериализации JSON в структуру Task.
func (t *Task) UnmarshalJSON(data []byte) error {
	var aux map[string]interface{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if id, ok := aux["id"].(float64); ok {
		t.ID = int64(id)
	}
	if text, ok := aux["text"].(string); ok {
		t.Text = text
	}
	if status, ok := aux["status"].(float64); ok {
		t.Status = int(status)
	}

	if createdDate, ok := aux["createdDate"].(string); ok {
		t.CreatedDate, _ = time.Parse("2006-01-02", createdDate)
	}
	if expectedDate, ok := aux["expectedDate"].(string); ok {
		t.ExpectedDate, _ = time.Parse("2006-01-02", expectedDate)
	}

	return nil
}
