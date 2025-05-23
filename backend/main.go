package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/rs/cors"
)

func main() {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{"https://imagemorph-abiola.netlify.app", "http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	http.HandleFunc("/api/convert", convertHandler)
	http.HandleFunc("/api/metadata", fetchVideoHandler)
	http.HandleFunc("/api/download", downloadHandler)
	http.HandleFunc("/api/removebg", removeBgHandler)
	
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))) 

	handler := corsOptions.Handler(http.DefaultServeMux)
	port := ":8080"

	fmt.Println("server is running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, handler))
}
