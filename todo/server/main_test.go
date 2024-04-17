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
	assert.NoError(t, err)
	defer mockDB.Close()

	db.DB = mockDB

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", handlers.GetTasks)
	mux.HandleFunc("/api/tasks/create", handlers.CreateTask)
	mux.HandleFunc("/api/tasks/update", handlers.UpdateTask)
	mux.HandleFunc("/api/tasks/delete", handlers.DeleteTask)

	srv := httptest.NewServer(mux)
	defer srv.Close()

	task := db.Task{
		Text:   "Test Task",
		Status: db.StatusInProgress,
	}
	taskJSON, _ := json.Marshal(task)

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.Text, sqlmock.AnyArg(), sqlmock.AnyArg(), task.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		srv.URL+"/api/tasks/create",
		bytes.NewBuffer(taskJSON),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdTask db.Task
	json.NewDecoder(resp.Body).Decode(&createdTask)
	assert.Equal(t, task.Text, createdTask.Text)
	assert.Equal(t, task.Status, createdTask.Status)

	fixedTime := time.Date(2023, time.April, 4, 11, 15, 0, 0, time.UTC)
	task1 := db.Task{ID: 1, Text: "Task 1", CreatedDate: fixedTime, ExpectedDate: fixedTime, Status: db.StatusInProgress}
	task2 := db.Task{ID: 2, Text: "Task 2", CreatedDate: fixedTime, ExpectedDate: fixedTime, Status: db.StatusCompleted}
	expectedTasks := []db.Task{task1, task2}

	rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"}).
		AddRow(task1.ID, task1.Text, task1.CreatedDate, task1.ExpectedDate, task1.Status).
		AddRow(task2.ID, task2.Text, task2.CreatedDate, task2.ExpectedDate, task2.Status)
	mock.ExpectQuery("SELECT id, task_text, createdDate, expectedDate, status FROM tasks").WillReturnRows(rows)

	req, err = http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualTasks []db.Task
	json.NewDecoder(resp.Body).Decode(&actualTasks)
	assert.Equal(t, expectedTasks, actualTasks)

	updatedTask := db.Task{
		ID:           1,
		Text:         "Updated Task",
		Status:       db.StatusCompleted,
		CreatedDate:  time.Now().UTC(),
		ExpectedDate: time.Now().UTC().AddDate(0, 0, 1),
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	mock.ExpectExec("UPDATE tasks SET task_text = \\$1, createdDate = \\$2, "+
		"expectedDate = \\$3, status = \\$4 WHERE id = \\$5").
		WithArgs(updatedTask.Text, updatedTask.CreatedDate, updatedTask.ExpectedDate, updatedTask.Status, updatedTask.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		srv.URL+"/api/tasks/update?id=1",
		bytes.NewBuffer(updatedTaskJSON),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ = http.NewRequestWithContext(context.Background(), http.MethodDelete, srv.URL+"/api/tasks/delete?id=1", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NoError(t, mock.ExpectationsWereMet())
}
