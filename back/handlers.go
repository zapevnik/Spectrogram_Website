package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
)

var (
	fileStore      = make(map[string]string) // Ключ: ID файла, Значение: путь к файлу
	fileStoreMutex = sync.Mutex{}            // Защита от одновременного доступа
)

// Отвечаем на запрос с фронтенда, отправляя изображение спектрограммы
func handleGetSpectrogram(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос получен на /getSpec")
	fileID := r.URL.Query().Get("fileid")
	if fileID == "" {
		http.Error(w, "ID файла не указан", http.StatusBadRequest)
		return
	}

	fileStoreMutex.Lock()
	fileName, exists := fileStore[fileID]
	fileStoreMutex.Unlock()

	if !exists {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		log.Printf("Файл с ID %s не найден", fileID)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Спектрограмма не найдена", http.StatusNotFound)
		log.Println("Ошибка открытия спектрограммы:", err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Ошибка отправки спектрограммы:", err)
	}
	file.Close()

	// Удаляем файл после использования
	defer func() {
		if err := os.Remove(fileName); err != nil {
			log.Printf("Ошибка при удалении файла %s: %v", fileName, err)
		}
	}()
}

// Скачиваем спектрограмму, полученную с Python-сервера
func handleUploadSpectrogram(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос получен на /uploadSpec")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := r.URL.Query().Get("fileid")
	if fileID == "" {
		http.Error(w, "Missing fileid parameter", http.StatusBadRequest)
		return
	}
	fileName := "./uploads/spectrogram_" + fileID + ".png"
	dst, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		log.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Println("Error saving file:", err)
		return
	}
	fileStoreMutex.Lock()
	fileStore[fileID] = fileName
	fileStoreMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

// Скачиваем аудиофайл, пришедший с фронтенда, и перенаправляем его на Python-сервер
func handleUploadFlac(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос получен на /uploadFlac")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	colormap := r.URL.Query().Get("colormap")
	if colormap == "" {
		log.Println("Missing colormap parameter")
		return
	}
	//Получаем файл и сохраняем в uploads
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		log.Println("Ошибка получения файла:", err)
		return
	}
	defer file.Close()

	fileID := uuid.New().String()
	fileName := fmt.Sprintf("./uploads/%s_%s", fileID, handler.Filename)

	dst, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		log.Println("Ошибка создания файла:", err)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		log.Println("Ошибка записи файла:", err)
		return
	}

	// Инициализируем ID файла
	fileStoreMutex.Lock()
	fileStore[fileID] = ""
	fileStoreMutex.Unlock()

	_, err = dst.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("Ошибка сброса указателя файла:", err)
		http.Error(w, "Ошибка обработки файла", http.StatusInternalServerError)
		return
	}

	// Отправляем файл
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", handler.Filename)
	if err != nil {
		log.Println("Ошибка создания form file:", err)
		http.Error(w, "Ошибка отправки файла", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(part, dst)
	if err != nil {
		log.Println("Ошибка копирования файла в форму:", err)
		http.Error(w, "Ошибка отправки файла", http.StatusInternalServerError)
		return
	}
	dst.Close()

	_ = writer.WriteField("colormap", colormap)

	writer.Close()

	url := fmt.Sprintf("http://localhost:5000/upload?fileid=%s&colormap=%s", fileID, colormap)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println("Ошибка создания запроса:", err)
		http.Error(w, "Ошибка отправки файла", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка отправки запроса:", err)
		http.Error(w, "Ошибка отправки файла", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"fileID": fileID}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	// Удаляем файл после использования
	defer func() {
		if err := os.Remove(fileName); err != nil {
			log.Printf("Ошибка при удалении файла %s: %v", fileName, err)
		}
	}()
}
