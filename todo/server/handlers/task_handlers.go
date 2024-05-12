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
	json.NewEncoder(w).Encode(tasks)
}

// Обработчик для создания новой задачи.
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Error decoding task: %v", err)
		http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received task data: %+v", task)

	task.Text = strings.TrimSpace(task.Text)
	if task.Text == "" {
		log.Println("Task text is empty")
		http.Error(w, `{"message": "Task text cannot be empty"}`, http.StatusBadRequest)
		return
	}

	if len(task.Text) > 255 {
		log.Println("Task text is too long")
		http.Error(w, `{"message": "Task text cannot exceed 255 characters"}`, http.StatusBadRequest)
		return
	}

	if task.ExpectedDate.Before(task.CreatedDate) {
		http.Error(w, "Expected date cannot be earlier than created date", http.StatusBadRequest)
		return
	}

	if task.Status != db.StatusInProgress &&
		task.Status != db.StatusCompleted &&
		task.Status != db.StatusTesting &&
		task.Status != db.StatusReturned {
		log.Printf("Invalid task status: %d", task.Status)
		http.Error(w, `{"message": "Invalid task status"}`, http.StatusBadRequest)
		return
	}

	err = db.CreateTask(task)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		http.Error(w, `{"message": "Failed to create task"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Task created successfully: %+v", task)

	w.WriteHeader(http.StatusCreated)
	jsonData, err := task.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// Обработчик для обновления существующей задачи.
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task db.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Error decoding task: "+err.Error(), http.StatusBadRequest)
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
	jsonData, err := task.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// Обработчик для удаления задачи.
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
