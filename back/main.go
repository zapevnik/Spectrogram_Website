package main

import (
	"fmt"
	"net/http"
	"os"
)

// Политика для общения с фронтендом
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/uploadFlac", handleUploadFlac)
	http.HandleFunc("/uploadSpec", handleUploadSpectrogram)
	http.HandleFunc("/getSpec", handleGetSpectrogram)

	fmt.Println("Сервер запущен на http://localhost:8080")

	err := http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux))
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

}
