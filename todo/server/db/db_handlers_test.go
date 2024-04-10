package db

import (
	"database/sql/driver"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Инициализация тестовой базы данных.
	InitDB()
	// Выполнение тестов.
	exitCode := m.Run()
	// Закрытие базы данных перед выходом.
	CloseDB()
	os.Exit(exitCode)
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestGetAllTasks(t *testing.T) {
	// Подготовка тестовых данных.
	task1 := Task{ID: 1, Text: "Task 1", Date: time.Now(), Status: "в процессе"}
	task2 := Task{ID: 2, Text: "Task 2", Date: time.Now(), Status: "завершено"}

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
			expectedTasks: []Task{task1, task2},
		},
		{
			name:          "Фильтрация по статусу 'в процессе'",
			statusFilter:  "в процессе",
			sortOrder:     "",
			expectedTasks: []Task{task1},
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
			expectedTasks: []Task{task1, task2},
		},
		{
			name:          "Сортировка по убыванию даты",
			statusFilter:  "",
			sortOrder:     "desc",
			expectedTasks: []Task{task2, task1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Мок базы данных.
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Замена глобальной переменной DB на мок базы данных.
			DB = db

			// Ожидаемый запрос и результат.
			rows := sqlmock.NewRows([]string{"id", "task_text", "task_date", "status"})
			for _, task := range tc.expectedTasks {
				rows.AddRow(task.ID, task.Text, task.Date, task.Status)
			}
			mock.ExpectQuery("SELECT id, task_text, task_date, status FROM tasks").WillReturnRows(rows)

			// Выполнение тестируемой функции.
			tasks, err := GetAllTasks(tc.statusFilter, tc.sortOrder)

			// Проверка результатов.
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTasks, tasks)
		})
	}
}

func TestCreateTask(t *testing.T) {
	task := Task{
		Text:   "New Task",
		Status: "в процессе",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.Text, AnyTime{}, task.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = CreateTask(task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTask(t *testing.T) {
	task := Task{
		ID:     1,
		Text:   "Updated Task",
		Status: "завершено",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(task.Text, AnyTime{}, task.Status, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = UpdateTask(task)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTask(t *testing.T) {
	taskID := 1

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = DeleteTask(taskID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTaskWithEmptyText(t *testing.T) {
	task := Task{
		Text:   "",
		Status: "в процессе",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.Text, AnyTime{}, task.Status).
		WillReturnError(errors.New("empty task text"))

	err = CreateTask(task)

	assert.Error(t, err)
	assert.EqualError(t, err, "empty task text")
}

func TestUpdateTaskWithNonExistentID(t *testing.T) {
	task := Task{
		ID:     999, // Несуществующий ID задачи.
		Text:   "Updated Task",
		Status: "завершено",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(task.Text, AnyTime{}, task.Status, task.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = UpdateTask(task)

	assert.Error(t, err)
	assert.EqualError(t, err, "task not found")
}
