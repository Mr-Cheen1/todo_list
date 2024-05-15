package db

import (
	"database/sql"
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

// Вспомогательная структура для сериализации Task.
type TaskDTO struct {
	ID           int64  `json:"id"`
	Text         string `json:"text"`
	CreatedDate  string `json:"createdDate"`
	ExpectedDate string `json:"expectedDate"`
	Status       int    `json:"status"`
}

// Метод для преобразования Task в TaskDTO.
func (t *Task) ToDTO() TaskDTO {
	return TaskDTO{
		ID:           t.ID,
		Text:         t.Text,
		CreatedDate:  t.CreatedDate.Format("2006-01-02"),
		ExpectedDate: t.ExpectedDate.Format("2006-01-02"),
		Status:       t.Status,
	}
}

// Метод для преобразования TaskDTO в Task.
func (dto *TaskDTO) ToTask() (Task, error) {
	createdDate, err := time.Parse("2006-01-02", dto.CreatedDate)
	if err != nil {
		return Task{}, err
	}
	expectedDate, err := time.Parse("2006-01-02", dto.ExpectedDate)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:           dto.ID,
		Text:         dto.Text,
		CreatedDate:  createdDate,
		ExpectedDate: expectedDate,
		Status:       dto.Status,
	}, nil
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
			// Если sortField не находитя в белом списке, возвращаем ошибку.
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

// Функция CreateTask создает новую задачу в базе данных и возвращает её ID.
func CreateTask(task Task) (int64, error) {
	query := "INSERT INTO tasks (task_text, createdDate, expectedDate, status) VALUES ($1, $2, $3, $4) RETURNING id"

	createdDateStr := task.CreatedDate.Format("2006-01-02")
	expectedDateStr := task.ExpectedDate.Format("2006-01-02")
	var id int64
	err := DB.QueryRow(query, task.Text, createdDateStr, expectedDateStr, task.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
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
