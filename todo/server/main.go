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

	// os.Setenv("DB_HOST", "localhost")
	// os.Setenv("DB_PORT", "8080")
	// os.Setenv("DB_USER", "postgres")
	// os.Setenv("DB_PASSWORD", "4217")
	// os.Setenv("DB_NAME", "todo_db")

	// fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	// fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	// fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	// fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	// fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

	db.InitDB()
	defer db.CloseDB()

	// Создание экземпляра сервера.
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", address, port),
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Регистрация обработчиков.
	// Для локального запуска.
	// Поменяйте на http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/", http.FileServer(http.Dir("/app/static")))
	http.HandleFunc("/api/tasks", handlers.GetTasks)
	http.HandleFunc("/api/tasks/create", handlers.CreateTask)
	http.HandleFunc("/api/tasks/update", handlers.UpdateTask)
	http.HandleFunc("/api/tasks/delete", handlers.DeleteTask)

	// Запуск сервера в отдельной горутине.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Printf("Server listening on %s:%s", address, port)

	// Ожидание сигнала для graceful shutdown.
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
