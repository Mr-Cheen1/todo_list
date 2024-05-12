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

// Тест для обработчика GetTasks.
func TestGetTasks(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer mockDB.Close()

	db.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "text", "createdDate", "expectedDate", "status"}).
		AddRow(1, "Test Task", time.Now(), time.Now().Add(24*time.Hour), db.StatusInProgress)
	mock.ExpectQuery("^SELECT (.+) FROM tasks").WillReturnRows(rows)

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"/tasks?status=active&sort=asc&sortField=createdDate",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"text":"Test Task","status":1,"createdDate":"...","expectedDate":"..."}]`
	if !strings.Contains(rr.Body.String(), "Test Task") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// Тест для обработчика CreateTask.
func TestCreateTask(t *testing.T) {
	taskText := "New Task"
	taskStatus := db.StatusInProgress
	createdDate := time.Now().UTC().Truncate(24 * time.Hour)
	expectedDate := createdDate.AddDate(0, 0, 1)
	taskJSON := fmt.Sprintf(`{"text":"%s","status":%d,"createdDate":"%s","expectedDate":"%s"}`,
		taskText, taskStatus, createdDate.Format("2006-01-02"), expectedDate.Format("2006-01-02"))

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(taskText, createdDate.Format("2006-01-02"), expectedDate.Format("2006-01-02"), taskStatus).
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
	assert.Equal(t, createdDate.Format("2006-01-02"), createdTask.CreatedDate.Format("2006-01-02"))
	assert.Equal(t, expectedDate.Format("2006-01-02"), createdTask.ExpectedDate.Format("2006-01-02"))

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест для обработчика UpdateTask.
func TestUpdateTask(t *testing.T) {
	fixedTime := time.Date(2023, time.April, 4, 0, 0, 0, 0, time.UTC)
	expectedTime := fixedTime.Add(48 * time.Hour)

	taskToUpdate := db.Task{
		ID:           1,
		Text:         "Updated Task",
		CreatedDate:  fixedTime,
		ExpectedDate: expectedTime,
		Status:       db.StatusInProgress,
	}

	taskJSON := fmt.Sprintf(`{"id":%d,"text":"%s","status":%d,"createdDate":"%s","expectedDate":"%s"}`,
		taskToUpdate.ID, taskToUpdate.Text, taskToUpdate.Status,
		taskToUpdate.CreatedDate.Format("2006-01-02"), taskToUpdate.ExpectedDate.Format("2006-01-02"))

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(`UPDATE tasks SET task_text = \$1, createdDate = \$2, 
		expectedDate = \$3, status = \$4 WHERE id = \$5`).
		WithArgs(taskToUpdate.Text, taskToUpdate.CreatedDate.Format("2006-01-02"),
			taskToUpdate.ExpectedDate.Format("2006-01-02"), taskToUpdate.Status, taskToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	db.DB = mockDB

	req, err := http.NewRequestWithContext(context.Background(), "PUT",
		"/api/tasks/update?id=1", strings.NewReader(taskJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateTask)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест для обработчика DeleteTask.
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
