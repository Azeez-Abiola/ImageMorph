package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"log"
	"github.com/joho/godotenv"
)

type RemovedBGImage struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

func removeBgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
 
  	err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found.")
		}
		
	apiKey := os.Getenv("REMOVEBG_API_KEY")
	if apiKey == "" {
		http.Error(w, "RemoveBG API key not configured", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, _ := writer.CreateFormFile("image_file", header.Filename)
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", &requestBody)
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to process image with RemoveBG", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("RemoveBG error: %s", string(body)), http.StatusBadRequest)
		return
	}

	uploadDir := "uploads/"
	os.MkdirAll(uploadDir, os.ModePerm)
	outputPath := filepath.Join(uploadDir, "bg_removed_"+header.Filename)
	outFile, _ := os.Create(outputPath)
	defer outFile.Close()

	io.Copy(outFile, resp.Body)
	fileInfo, _ := outFile.Stat()

	response := struct {
		Message string         `json:"message"`
		Image   RemovedBGImage `json:"image"`
	}{
		Message: "Background removed successfully",
		Image: RemovedBGImage{
			URL:  "/" + outputPath,
			Size: fileInfo.Size(),
		},
	}


	go func() {
		time.Sleep(23 * time.Hour)
		os.Remove(outputPath)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
