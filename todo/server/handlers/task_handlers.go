package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mr-Cheen1/todo_list/server/db"
)

// Обработчик для получения списка задач.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	sortOrder := r.URL.Query().Get("sort")
	sortField := r.URL.Query().Get("sortField")
	tasks, err := db.GetAllTasks(statusFilter, sortOrder, sortField)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskDTOs := make([]db.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, task.ToDTO())
	}

	json.NewEncoder(w).Encode(taskDTOs)
}

// Обработчик для создания новой задачи.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO db.TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := taskDTO.ToTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.Text = strings.TrimSpace(task.Text)
	if task.Text == "" {
		http.Error(w, "Task text cannot be empty", http.StatusBadRequest)
		return
	}

	if task.ExpectedDate.Before(task.CreatedDate) {
		http.Error(w, "Expected date cannot be earlier than created date", http.StatusBadRequest)
		return
	}

	if len(task.Text) > 255 {
		log.Println("Task text is too long")
		http.Error(w, "Task text cannot exceed 255 characters", http.StatusBadRequest)
		return
	}

	if task.Status != db.StatusInProgress &&
		task.Status != db.StatusCompleted &&
		task.Status != db.StatusTesting &&
		task.Status != db.StatusReturned {
		http.Error(w, "Incorrect task status", http.StatusBadRequest)
		return
	}

	if task.CreatedDate.IsZero() {
		http.Error(w, "Task created date is required", http.StatusBadRequest)
		return
	}

	if task.ExpectedDate.IsZero() {
		http.Error(w, "Task expected date is required", http.StatusBadRequest)
		return
	}

	id, err := db.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task.ToDTO())
}

// Обработчик для обновления существующей задачи.
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var taskDTO db.TaskDTO
	err = json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		http.Error(w, "Error decoding task: "+err.Error(), http.StatusBadRequest)
		return
	}

	task, err := taskDTO.ToTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.ID = int64(id)

	log.Printf("Updating task: %+v", task)

	task.Text = strings.TrimSpace(task.Text)
	if task.Text == "" {
		http.Error(w, "Task text cannot be empty", http.StatusBadRequest)
		return
	}

	if task.ExpectedDate.Before(task.CreatedDate) {
		http.Error(w, "Expected date cannot be earlier than created date", http.StatusBadRequest)
		return
	}

	if len(task.Text) > 255 {
		log.Println("Task text is too long")
		http.Error(w, "Task text cannot exceed 255 characters", http.StatusBadRequest)
		return
	}

	if task.Status != db.StatusInProgress &&
		task.Status != db.StatusCompleted &&
		task.Status != db.StatusTesting &&
		task.Status != db.StatusReturned {
		http.Error(w, "Incorrect task status", http.StatusBadRequest)
		return
	}

	if task.CreatedDate.IsZero() {
		http.Error(w, "Task created date is required", http.StatusBadRequest)
		return
	}

	if task.ExpectedDate.IsZero() {
		http.Error(w, "Task expected date is required", http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task updated successfully: %+v", task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task.ToDTO())
}

// Обработчик для удаления задачи.
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
