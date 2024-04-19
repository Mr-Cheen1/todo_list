package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mr-Cheen1/todo_list/server/db"
	"github.com/Mr-Cheen1/todo_list/server/handlers"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	// Настройка моков для GetAllTasks.
	rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"}).
		AddRow(1, "Test Task", time.Now(), time.Now().Add(24*time.Hour), db.StatusInProgress)

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
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs("New Task", sqlmock.AnyArg(), sqlmock.AnyArg(), db.StatusInProgress).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Тестирование CreateTask.
	newTask := db.Task{
		Text:         "New Task",
		CreatedDate:  time.Now(),
		ExpectedDate: time.Now().Add(24 * time.Hour),
		Status:       db.StatusInProgress,
	}
	body, _ := json.Marshal(newTask)
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
	json.NewDecoder(w.Body).Decode(&createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "New Task", createdTask.Text)

	// Настройка моков для UpdateTask.
	mock.ExpectExec("UPDATE tasks SET").
		WithArgs("Updated Task", sqlmock.AnyArg(), sqlmock.AnyArg(), db.StatusCompleted, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Тестирование UpdateTask.
	updateTask := db.Task{
		ID:           1,
		Text:         "Updated Task",
		CreatedDate:  time.Now(),
		ExpectedDate: time.Now().Add(48 * time.Hour),
		Status:       db.StatusCompleted,
	}
	body, _ = json.Marshal(updateTask)
	req, _ = http.NewRequestWithContext(
		context.Background(),
		"PUT",
		server.URL+"/tasks/update?id=1",
		bytes.NewBuffer(body),
	)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Настройка моков для DeleteTask.
	mock.ExpectExec("DELETE FROM tasks WHERE").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	// Тестирование DeleteTask.
	req, _ = http.NewRequestWithContext(
		context.Background(),
		"DELETE",
		server.URL+"/tasks/delete?id=1",
		nil,
	)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
