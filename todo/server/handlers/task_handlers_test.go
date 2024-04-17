package handlers

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
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	fixedTime := time.Date(2023, time.April, 4, 11, 15, 0, 0, time.UTC)
	task1 := db.Task{ID: 1, Text: "Task 1", CreatedDate: fixedTime, Status: db.StatusInProgress}
	task2 := db.Task{ID: 2, Text: "Task 2", CreatedDate: fixedTime, Status: db.StatusCompleted}
	expectedTasks := []db.Task{task1, task2}

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"}).
		AddRow(task1.ID, task1.Text, task1.CreatedDate, task1.ExpectedDate, task1.Status).
		AddRow(task2.ID, task2.Text, task2.CreatedDate, task2.ExpectedDate, task2.Status)
	mock.ExpectQuery("SELECT id, task_text, createdDate, expectedDate, status FROM tasks").WillReturnRows(rows)

	db.DB = mockDB

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetTasks)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var actualTasks []db.Task
	err = json.Unmarshal(rr.Body.Bytes(), &actualTasks)
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, actualTasks)
}

func TestCreateTask(t *testing.T) {
	taskText := "New Task"
	taskStatus := db.StatusInProgress
	createdDate := time.Now().UTC()
	expectedDate := createdDate.AddDate(0, 0, 1)
	taskJSON := fmt.Sprintf(`{"text":"%s","status":%d,"createdDate":"%s","expectedDate":"%s"}`,
		taskText, taskStatus, createdDate.Format(time.RFC3339), expectedDate.Format(time.RFC3339))

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(taskText, sqlmock.AnyArg(), sqlmock.AnyArg(), taskStatus).
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
	assert.WithinDuration(t, createdDate, createdTask.CreatedDate, time.Second)
	assert.WithinDuration(t, expectedDate, createdTask.ExpectedDate, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	fixedTime := time.Date(2023, time.April, 4, 11, 15, 0, 0, time.UTC)
	expectedTime := fixedTime.Add(48 * time.Hour)

	taskToUpdate := db.Task{
		ID:           1,
		Text:         "Updated Task",
		CreatedDate:  fixedTime,
		ExpectedDate: expectedTime,
		Status:       db.StatusInProgress,
	}

	taskJSON, err := json.Marshal(taskToUpdate)
	if err != nil {
		t.Fatalf("Error marshaling task: %v", err)
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`
    UPDATE tasks 
    SET task_text = \$1, createdDate = \$2, expectedDate = \$3, status = \$4 
    WHERE id = \$5
`).WithArgs(taskToUpdate.Text, fixedTime, expectedTime, taskToUpdate.Status, taskToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	db.DB = mockDB

	req, err := http.NewRequestWithContext(
		context.Background(),
		"PUT",
		"/api/tasks/update?id=1",
		bytes.NewBuffer(taskJSON),
	)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
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
