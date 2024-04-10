package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mr-Cheen1/todo_list/server/db"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	sortOrder := r.URL.Query().Get("sort")
	tasks, err := db.GetAllTasks(statusFilter, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Incorrect data format", http.StatusBadRequest)
		return
	}

	// Валидация полей задачи.
	if task.Text == "" {
		http.Error(w, "The task text cannot be empty", http.StatusBadRequest)
		return
	}

	if task.Status != "в процессе" && task.Status != "завершено" {
		http.Error(w, "Incorrect task status", http.StatusBadRequest)
		return
	}

	err = db.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var task db.Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = id
	err := db.UpdateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
