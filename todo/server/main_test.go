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

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper() // Add this line to mark the function as a test helper
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	db.DB = mockDB
	return mock, func() { mockDB.Close() }
}

func setupServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/get", handlers.GetTasks)
	mux.HandleFunc("/tasks/create", handlers.CreateTask)
	mux.HandleFunc("/tasks/update", handlers.UpdateTask)
	mux.HandleFunc("/tasks/delete", handlers.DeleteTask)
	return httptest.NewServer(mux)
}

func TestGetTasks(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	fixedTime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"}).
		AddRow(1, "Test Task", fixedTime, fixedTime.Add(24*time.Hour), db.StatusInProgress)
	mock.ExpectQuery("^SELECT (.+) FROM tasks$").WillReturnRows(rows)

	server := setupServer()
	defer server.Close()

	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL+"/tasks/get", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	w := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var taskDTOs []db.TaskDTO
	err = json.NewDecoder(w.Body).Decode(&taskDTOs)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	assert.Equal(t, 1, len(taskDTOs))
	assert.Equal(t, "Test Task", taskDTOs[0].Text)
}

func TestCreateTask(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	createdDate := time.Now().Truncate(24 * time.Hour)
	expectedDate := createdDate.AddDate(0, 0, 1)
	mock.ExpectQuery(
		"INSERT INTO tasks \\(task_text, createdDate, expectedDate, status\\) "+
			"VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id",
	).
		WithArgs(
			"New Task",
			createdDate.Format("2006-01-02"),
			expectedDate.Format("2006-01-02"),
			db.StatusInProgress,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	server := setupServer()
	defer server.Close()

	newTaskDTO := db.TaskDTO{
		Text:         "New Task",
		Status:       db.StatusInProgress,
		CreatedDate:  createdDate.Format("2006-01-02"),
		ExpectedDate: expectedDate.Format("2006-01-02"),
	}

	body, err := json.Marshal(newTaskDTO)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		server.URL+"/tasks/create",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	w := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(w, req)

	t.Logf("Response body: %s", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask db.TaskDTO
	err = json.NewDecoder(w.Body).Decode(&createdTask)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	assert.Equal(t, "New Task", createdTask.Text)
	assert.Equal(t, db.StatusInProgress, createdTask.Status)
	assert.Equal(t, createdDate.Format("2006-01-02"), createdTask.CreatedDate)
	assert.Equal(t, expectedDate.Format("2006-01-02"), createdTask.ExpectedDate)
}

func TestUpdateTask(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	fixedTime := time.Now()
	expectedTime := fixedTime.Add(48 * time.Hour)

	taskToUpdate := db.TaskDTO{
		ID:           1,
		Text:         "Updated Task",
		CreatedDate:  fixedTime.Format("2006-01-02"),
		ExpectedDate: expectedTime.Format("2006-01-02"),
		Status:       db.StatusInProgress,
	}

	mock.ExpectExec(`UPDATE tasks SET task_text = \$1, createdDate = \$2, `+
		`expectedDate = \$3, status = \$4 WHERE id = \$5`).
		WithArgs(
			taskToUpdate.Text,
			taskToUpdate.CreatedDate,
			taskToUpdate.ExpectedDate,
			taskToUpdate.Status,
			taskToUpdate.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	server := setupServer()
	defer server.Close()

	taskJSON := fmt.Sprintf(`{"id":%d,"text":"%s","status":%d,"createdDate":"%s","expectedDate":"%s"}`,
		taskToUpdate.ID, taskToUpdate.Text, taskToUpdate.Status,
		taskToUpdate.CreatedDate, taskToUpdate.ExpectedDate)

	req, err := http.NewRequestWithContext(
		context.Background(),
		"PUT",
		server.URL+"/tasks/update?id=1",
		strings.NewReader(taskJSON),
	)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTask(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	taskIDToDelete := 1
	mock.ExpectExec(`^DELETE FROM tasks WHERE id = \$1$`).
		WithArgs(taskIDToDelete).
		WillReturnResult(sqlmock.NewResult(0, 1))

	server := setupServer()
	defer server.Close()

	req, err := http.NewRequestWithContext(
		context.Background(),
		"DELETE",
		fmt.Sprintf("%s/tasks/delete?id=%d", server.URL, taskIDToDelete),
		nil,
	)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	w := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(w, req)

	t.Logf("Response body: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}
