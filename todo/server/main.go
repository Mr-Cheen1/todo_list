package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mr-Cheen1/todo_list/server/db"
	"github.com/Mr-Cheen1/todo_list/server/handlers"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server/main.go <address> <port>")
		os.Exit(1)
	}

	address := os.Args[1]
	port := os.Args[2]

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.CloseDB()

	// Создание экземпляра сервера.
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", address, port),
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
	}

	http.Handle("/", http.FileServer(http.Dir("/app/static")))
	http.HandleFunc("/api/tasks", handlers.GetTasks)
	http.HandleFunc("/api/tasks/create", handlers.CreateTask)
	http.HandleFunc("/api/tasks/update", handlers.UpdateTask)
	http.HandleFunc("/api/tasks/delete", handlers.DeleteTask)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Printf("Server listening on %s:%s", address, port)

	if err := gracefulShutdown(srv); err != nil {
		log.Println("Failed to gracefully shutdown:", err)
		return
	}
}

func gracefulShutdown(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
		return err
	}
	log.Println("Server exiting")
	return nil
}
