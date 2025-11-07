package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// HomeHandler — отдаёт index.html
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// UploadHandler — обрабатывает загрузку файла
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму (до 10 МБ)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		if err == http.ErrContentLength {
			http.Error(w, "Файл слишком большой (макс. 10 МБ)", http.StatusRequestEntityTooLarge)
			return
		}
		log.Printf("Ошибка парсинга формы: %v", err)
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	// Получаем файл
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		http.Error(w, "Ошибка получения файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем содержимое
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка чтения файла: %v", err)
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	input := string(data)

	// Конвертируем
	result, err := service.Convert(input)
	if err != nil {
		log.Printf("Ошибка конвертации: %v", err)
		http.Error(w, "Ошибка конвертации данных", http.StatusInternalServerError)
		return
	}

	// Генерим имя выходного файла
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		ext = ".txt"
	}
	outputFilename := time.Now().UTC().Format("2006-01-02-15-04-05") + ext

	// Записываем в локальный файл
	err = os.WriteFile(outputFilename, []byte(result), 0644)
	if err != nil {
		log.Printf("Ошибка записи файла: %v", err)
		http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}
