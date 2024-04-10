package main_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mr-Cheen1/todo_list/server/db"
	"github.com/Mr-Cheen1/todo_list/server/handlers"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Создание мока базы данных.
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Замена глобальной переменной DB на мок.
	db.DB = mockDB

	// Создание тестового сервера.
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.HandleFunc("/api/tasks", handlers.GetTasks)
	router.HandleFunc("/api/tasks/create", handlers.CreateTask)
	router.HandleFunc("/api/tasks/update", handlers.UpdateTask)
	router.HandleFunc("/api/tasks/delete", handlers.DeleteTask)
	srv := httptest.NewServer(router)
	defer srv.Close()

	// Тест создания задачи.
	taskText := "Test Task"
	taskStatus := "в процессе"
	taskJSON := `{"text":"` + taskText + `","status":"` + taskStatus + `"}`

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(taskText, sqlmock.AnyArg(), taskStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		srv.URL+"/api/tasks/create",
		strings.NewReader(taskJSON),
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var createdTask db.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, taskText, createdTask.Text)
	assert.Equal(t, taskStatus, createdTask.Status)

	// Тест обновления задачи.
	updatedTaskText := "Updated Task"
	updatedTaskStatus := "выполнено"
	updatedTaskJSON := `{"id":` + strconv.Itoa(createdTask.ID) +
		`,"text":"` + updatedTaskText +
		`","status":"` + updatedTaskStatus +
		`","date":"` + createdTask.Date.Format("2006-01-02T15:04:05Z07:00") + `"}`

	mock.ExpectExec("UPDATE tasks").
		WithArgs(updatedTaskText, createdTask.Date, updatedTaskStatus, createdTask.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		srv.URL+"/api/tasks/update?id="+strconv.Itoa(createdTask.ID),
		strings.NewReader(updatedTaskJSON),
	)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// Обновление ожидаемого значения.
	expectedTask := db.Task{
		ID:     createdTask.ID,
		Text:   updatedTaskText,
		Date:   createdTask.Date,
		Status: updatedTaskStatus,
	}

	// Тест получения списка задач.
	rows := sqlmock.NewRows([]string{"id", "task_text", "task_date", "status"}).
		AddRow(expectedTask.ID, expectedTask.Text, expectedTask.Date, expectedTask.Status)
	mock.ExpectQuery("SELECT id, task_text, task_date, status FROM tasks").WillReturnRows(rows)

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/api/tasks", nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var tasks []db.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, expectedTask, tasks[0])

	// Тест удаления задачи.
	mock.ExpectExec("DELETE FROM tasks").WithArgs(createdTask.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		srv.URL+"/api/tasks/delete?id="+strconv.Itoa(createdTask.ID),
		nil,
	)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// Проверка, что задача удалена.
	mock.ExpectQuery("SELECT id, task_text, task_date, status FROM tasks").
		WillReturnRows(sqlmock.NewRows([]string{"id", "task_text", "task_date", "status"}))

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/api/tasks", nil)
	assert.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	tasks = []db.Task{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 0)

	// Проверка вызовов моков.
	assert.NoError(t, mock.ExpectationsWereMet())
}
