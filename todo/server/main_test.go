package main

import (
	"bytes"
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
	"github.com/Mr-Cheen1/todo_list/server/handlers"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Создание мока базы данных.
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	// Настройка моков для GetAllTasks.
	fixedTime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"}).
		AddRow(1, "Test Task", fixedTime, fixedTime.Add(24*time.Hour), db.StatusInProgress)

	mock.ExpectQuery("^SELECT (.+) FROM tasks$").WillReturnRows(rows)

	// Создание тестового сервера.
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/get", handlers.GetTasks)
	mux.HandleFunc("/tasks/create", handlers.CreateTask)
	mux.HandleFunc("/tasks/update", handlers.UpdateTask)
	mux.HandleFunc("/tasks/delete", handlers.DeleteTask)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Тестирование GetTasks.
	req, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL+"/tasks/get", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tasks []db.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "Test Task", tasks[0].Text)

	// Настройка моков для CreateTask.
	createdDate := time.Now().Truncate(24 * time.Hour)
	expectedDate := createdDate.AddDate(0, 0, 1)
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs("New Task", createdDate.Format("2006-01-02"), expectedDate.Format("2006-01-02"), db.StatusInProgress).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Тестирование CreateTask.
	newTask := db.Task{
		Text:   "New Task",
		Status: db.StatusInProgress,
	}
	newTask.CreatedDate = time.Now().Truncate(24 * time.Hour)
	newTask.ExpectedDate = newTask.CreatedDate.AddDate(0, 0, 1)

	// Создаем новый объект Task с сериализованными датами.
	newTaskJSON := struct {
		Text         string `json:"text"`
		Status       int    `json:"status"`
		CreatedDate  string `json:"createdDate"`
		ExpectedDate string `json:"expectedDate"`
	}{
		Text:         newTask.Text,
		Status:       newTask.Status,
		CreatedDate:  newTask.CreatedDate.Format("2006-01-02"),
		ExpectedDate: newTask.ExpectedDate.Format("2006-01-02"),
	}

	body, _ := json.Marshal(newTaskJSON)
	req, _ = http.NewRequestWithContext(
		context.Background(),
		"POST",
		server.URL+"/tasks/create",
		bytes.NewBuffer(body),
	)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask db.Task
	err = json.NewDecoder(w.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "New Task", createdTask.Text)
	assert.Equal(t, db.StatusInProgress, createdTask.Status)
	assert.Equal(t, newTask.CreatedDate.Format("2006-01-02"), createdTask.CreatedDate.Format("2006-01-02"))
	assert.Equal(t, newTask.ExpectedDate.Format("2006-01-02"), createdTask.ExpectedDate.Format("2006-01-02"))

	// Настройка моков для UpdateTask.
	fixedTime = time.Now()
	expectedTime := fixedTime.Add(48 * time.Hour)

	taskToUpdate := db.Task{
		ID:           1,
		Text:         "Updated Task",
		CreatedDate:  fixedTime,
		ExpectedDate: expectedTime,
		Status:       db.StatusInProgress,
	}

	mock.ExpectExec(`UPDATE tasks SET task_text = \$1, createdDate = \$2, `+
		`expectedDate = \$3, status = \$4 WHERE id = \$5`).
		WithArgs(taskToUpdate.Text, taskToUpdate.CreatedDate.Format("2006-01-02"),
			taskToUpdate.ExpectedDate.Format("2006-01-02"), taskToUpdate.Status, taskToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Тестирование UpdateTask.
	taskJSON := fmt.Sprintf(`{"id":%d,"text":"%s","status":%d,"createdDate":"%s","expectedDate":"%s"}`,
		taskToUpdate.ID, taskToUpdate.Text, taskToUpdate.Status,
		taskToUpdate.CreatedDate.Format("2006-01-02"), taskToUpdate.ExpectedDate.Format("2006-01-02"))

	req, _ = http.NewRequestWithContext(
		context.Background(),
		"PUT",
		server.URL+"/tasks/update?id=1",
		strings.NewReader(taskJSON),
	)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Настройка моков для DeleteTask.
	taskIDToDelete := 1
	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs(taskIDToDelete).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Тестирование DeleteTask.
	req, _ = http.NewRequestWithContext(
		context.Background(),
		"DELETE",
		fmt.Sprintf("%s/tasks/delete?id=%d", server.URL, taskIDToDelete),
		nil,
	)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка выполнения всех ожиданий моков.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
