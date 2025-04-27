package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)

type VideoRequest struct {
	URL string `json:"url"`
}

type VideoResponse struct {
	URL string `json:"url"`
	Source string `json:"source"`
	ID string `json:"id"`
	Author string `json:"author"`
	Title string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Duration int  `json:"duration"`
	Medias []struct {
		URL string `json:"url"`
		Quality string `json:"quality"`
		Width int `json:"width"`
		Height int `json:"height"`
		Ext string `json:"ext"`
	} `json:"medias"`
	Error bool `json:"error"`
}


func fetchVideoMetaData(videoURL string) (*VideoResponse, error) {
 		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found.")
		}

 apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RAPIDAPI_KEY not set")
	}

	payload := VideoRequest{URL: videoURL}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://social-download-all-in-one.p.rapidapi.com/v1/social/autolink",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", apiKey) // 🔒 Backend-only key
	req.Header.Add("x-rapidapi-host", "social-download-all-in-one.p.rapidapi.com")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result VideoResponse
	if json.Unmarshal(body, &result) != nil || result.Error {
		return &result, nil
	}

	return &result, nil
}

func fetchVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	videoData, err := fetchVideoMetaData(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videoData)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter required", http.StatusBadRequest)
		return
	}

	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		fileName = "video.mp4"
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body) 
}


