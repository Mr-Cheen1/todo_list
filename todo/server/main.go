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
	// Проверка аргументов командной строки.
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run server/main.go <address> <port>")
		os.Exit(1)
	}

	address := os.Args[1]
	port := os.Args[2]

	// Получение параметров подключения к базе данных из переменных окружения.
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Инициализация подключения к базе данных.
	db.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.CloseDB()

	// Создание экземпляра сервера.
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", address, port),
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Регистрация обработчиков маршрутов.
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

	// Ожидание сигнала завершения и корректное завершение работы сервера.
	if err := gracefulShutdown(srv); err != nil {
		log.Println("Failed to gracefully shutdown:", err)
		return
	}
}

func gracefulShutdown(srv *http.Server) error {
	// Ожидание сигнала завершения.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Установка таймаута для корректного завершения работы сервера.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Корректное завершение работы сервера.
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}
