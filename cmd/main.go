package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LanaAntonova/go-sprint6-http/server"
)

func main() {
	logger := log.New(os.Stderr, "MORSE: ", log.LstdFlags|log.Lshortfile)

	srv := server.NewServer(logger)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := srv.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
