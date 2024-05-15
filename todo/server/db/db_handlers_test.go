package db

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Структура AnyTime используется для сопоставления любых значений типа time.Time в тестах.
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// Тест для функции GetAllTasks.
func TestGetAllTasks(t *testing.T) {
	// Подготовка тестовых данных.
	fixedTime := time.Now()
	task1 := Task{ID: 1, Text: "Task 1", CreatedDate: fixedTime, ExpectedDate: fixedTime, Status: StatusInProgress}
	task2 := Task{ID: 2, Text: "Task 2", CreatedDate: fixedTime, ExpectedDate: fixedTime, Status: StatusCompleted}
	task3 := Task{ID: 3, Text: "Task 3", CreatedDate: fixedTime, ExpectedDate: fixedTime, Status: StatusInProgress}

	testCases := []struct {
		name          string
		statusFilter  string
		sortOrder     string
		expectedTasks []Task
	}{
		{
			name:          "Получить все задачи",
			statusFilter:  "",
			sortOrder:     "",
			expectedTasks: []Task{task1, task2, task3},
		},
		{
			name:          "Фильтрация по статусу 'в процессе'",
			statusFilter:  "в процессе",
			sortOrder:     "",
			expectedTasks: []Task{task1, task3},
		},
		{
			name:          "Фильтрация по статусу 'завершено'",
			statusFilter:  "завершено",
			sortOrder:     "",
			expectedTasks: []Task{task2},
		},
		{
			name:          "Сортировка по возрастанию даты",
			statusFilter:  "",
			sortOrder:     "asc",
			expectedTasks: []Task{task1, task2, task3},
		},
		{
			name:          "Сортировка по убыванию даты",
			statusFilter:  "",
			sortOrder:     "desc",
			expectedTasks: []Task{task3, task2, task1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			DB = db // Замена глобальной переменной DB на мок базы данных.

			// Настройка ожидаемого запроса и возвращаемых данных.
			rows := sqlmock.NewRows([]string{"id", "task_text", "createdDate", "expectedDate", "status"})
			for _, task := range tc.expectedTasks {
				rows.AddRow(task.ID, task.Text, task.CreatedDate, task.ExpectedDate, task.Status)
			}
			mock.ExpectQuery("SELECT id, task_text, createdDate, expectedDate, status FROM tasks").WillReturnRows(rows)

			// Вызов тестируемой функции.
			tasks, err := GetAllTasks(tc.statusFilter, tc.sortOrder, "")

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTasks, tasks)
		})
	}
}

// Тест для функции CreateTask.
func TestCreateTask(t *testing.T) {
	task := Task{
		Text:         "New Task",
		Status:       StatusInProgress,
		CreatedDate:  time.Now().Truncate(24 * time.Hour),
		ExpectedDate: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	createdDateStr := task.CreatedDate.Format("2006-01-02")
	expectedDateStr := task.ExpectedDate.Format("2006-01-02")

	// Настройка ожидаемого запроса и возвращаемого результата.
	mock.ExpectQuery("INSERT INTO tasks").
		WithArgs(task.Text, createdDateStr, expectedDateStr, task.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызов тестируемой функции.
	id, err := CreateTask(task)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест для функции UpdateTask.
func TestUpdateTask(t *testing.T) {
	task := Task{
		ID:           1,
		Text:         "Updated Task",
		Status:       StatusCompleted,
		CreatedDate:  time.Now().Truncate(24 * time.Hour),
		ExpectedDate: time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour),
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	createdDateStr := task.CreatedDate.Format("2006-01-02")
	expectedDateStr := task.ExpectedDate.Format("2006-01-02")

	// Настройка ожидаемого запроса и возвращаемого результата.
	mock.ExpectExec("UPDATE tasks SET task_text = \\$1, createdDate = \\$2, "+
		"expectedDate = \\$3, status = \\$4 WHERE id = \\$5").
		WithArgs(task.Text, createdDateStr, expectedDateStr, task.Status, task.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов тестируемой функции.
	err = UpdateTask(task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест для функции DeleteTask.
func TestDeleteTask(t *testing.T) {
	taskID := 1

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	// Настройка ожидаемого запроса и возвращаемого результата.
	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Вызов тестируемой функции.
	err = DeleteTask(taskID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskToDTO(t *testing.T) {
	task := Task{
		ID:           1,
		Text:         "Test task",
		CreatedDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		ExpectedDate: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
		Status:       StatusInProgress,
	}

	dto := task.ToDTO()

	if dto.ID != task.ID {
		t.Errorf("expected ID %d, got %d", task.ID, dto.ID)
	}
	if dto.Text != task.Text {
		t.Errorf("expected Text %s, got %s", task.Text, dto.Text)
	}
	if dto.CreatedDate != "2023-10-01" {
		t.Errorf("expected CreatedDate 2023-10-01, got %s", dto.CreatedDate)
	}
	if dto.ExpectedDate != "2023-10-10" {
		t.Errorf("expected ExpectedDate 2023-10-10, got %s", dto.ExpectedDate)
	}
	if dto.Status != task.Status {
		t.Errorf("expected Status %d, got %d", task.Status, dto.Status)
	}
}

func TestTaskDTOToTask(t *testing.T) {
	dto := TaskDTO{
		ID:           1,
		Text:         "Test task",
		CreatedDate:  "2023-10-01",
		ExpectedDate: "2023-10-10",
		Status:       StatusInProgress,
	}

	task, err := dto.ToTask()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCreatedDate := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	expectedExpectedDate := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)

	if task.ID != dto.ID {
		t.Errorf("expected ID %d, got %d", dto.ID, task.ID)
	}
	if task.Text != dto.Text {
		t.Errorf("expected Text %s, got %s", dto.Text, task.Text)
	}
	if !task.CreatedDate.Equal(expectedCreatedDate) {
		t.Errorf("expected CreatedDate %v, got %v", expectedCreatedDate, task.CreatedDate)
	}
	if !task.ExpectedDate.Equal(expectedExpectedDate) {
		t.Errorf("expected ExpectedDate %v, got %v", expectedExpectedDate, task.ExpectedDate)
	}
	if task.Status != dto.Status {
		t.Errorf("expected Status %d, got %d", dto.Status, task.Status)
	}
}
