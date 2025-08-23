package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Путь к frontend файлам
	frontendDir := "../../frontend"
	if len(os.Args) > 1 {
		frontendDir = os.Args[1]
	}

	// Абсолютный путь
	absPath, err := filepath.Abs(frontendDir)
	if err != nil {
		log.Fatal("Ошибка получения пути:", err)
	}

	// Проверяем существование директории
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatal("Директория frontend не найдена:", absPath)
	}

	// Создаем файловый сервер
	fileServer := http.FileServer(http.Dir(absPath))
	
	// Настраиваем маршруты с отключением кеша для разработки
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отключаем кеширование для .js и .css файлов
		if strings.HasSuffix(r.URL.Path, ".js") || strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		fileServer.ServeHTTP(w, r)
	}))
	
	port := "3000"
	if envPort := os.Getenv("FRONTEND_PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🌐 Frontend сервер запущен на порту %s\n", port)
	fmt.Printf("📂 Статические файлы: %s\n", absPath)
	fmt.Printf("🔗 Откройте: http://localhost:%s\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
