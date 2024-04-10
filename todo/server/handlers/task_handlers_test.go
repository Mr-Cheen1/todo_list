package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mr-Cheen1/todo_list/server/db"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	// Подготовка тестовых данных.
	fixedTime := time.Date(2023, time.April, 4, 11, 15, 0, 0, time.UTC)
	task1 := db.Task{ID: 1, Text: "Task 1", Date: fixedTime, Status: "в процессе"}
	task2 := db.Task{ID: 2, Text: "Task 2", Date: fixedTime, Status: "завершено"}
	expectedTasks := []db.Task{task1, task2}

	// Мок базы данных.
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Ожидаемый запрос и результат.
	rows := sqlmock.NewRows([]string{"id", "task_text", "task_date", "status"}).
		AddRow(task1.ID, task1.Text, task1.Date, task1.Status).
		AddRow(task2.ID, task2.Text, task2.Date, task2.Status)
	mock.ExpectQuery("SELECT id, task_text, task_date, status FROM tasks").WillReturnRows(rows)

	// Заменяем глобальную переменную DB на мок базы данных.
	db.DB = mockDB

	// Выполнение тестируемой функции.
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/api/tasks", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks)
	handler.ServeHTTP(rr, req)

	// Проверка результатов.
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualTasks []db.Task
	err = json.Unmarshal(rr.Body.Bytes(), &actualTasks)
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, actualTasks)
}

func TestCreateTask(t *testing.T) {
	taskText := "New Task"
	taskStatus := "в процессе"
	taskJSON := fmt.Sprintf(`{"text":"%s","status":"%s"}`, taskText, taskStatus)

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(taskText, sqlmock.AnyArg(), taskStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db.DB = mockDB

	req, err := http.NewRequestWithContext(context.Background(), "POST", "/api/tasks/create", strings.NewReader(taskJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdTask db.Task
	err = json.Unmarshal(rr.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, taskText, createdTask.Text)
	assert.Equal(t, taskStatus, createdTask.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	taskID := 1
	taskText := "Updated Task"
	taskStatus := "завершено"
	taskJSON := fmt.Sprintf(`{"text":"%s","status":"%s"}`, taskText, taskStatus)

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(taskText, sqlmock.AnyArg(), taskStatus, taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	db.DB = mockDB

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		fmt.Sprintf("/api/tasks/update?id=%d", taskID),
		strings.NewReader(taskJSON),
	)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTask(t *testing.T) {
	taskID := 1

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").WithArgs(taskID).WillReturnResult(sqlmock.NewResult(0, 1))

	db.DB = mockDB

	deleteURL := fmt.Sprintf("/api/tasks/delete?id=%d", taskID)
	req, err := http.NewRequestWithContext(context.Background(), "DELETE", deleteURL, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(DeleteTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}
