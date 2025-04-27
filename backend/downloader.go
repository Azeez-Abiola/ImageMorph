package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func fetchVideoMetaData(videoURL, apiKey string) (*VideoResponse, error) {
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

	req.Header.Add("x-rapidapi-key", apiKey)
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
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
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error streaming video: %v", err)
	}
}

