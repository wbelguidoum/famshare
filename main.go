package main

import (
	"log"
	"net/http"
)

// --- API Handlers (Unchanged) ---

func main() {
	mux := http.NewServeMux()

	// Static files and uploads
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Public API routes
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/logout", logoutHandler)

	// Protected API routes
	mux.Handle("/api/session/check", authMiddleware(http.HandlerFunc(checkSessionHandler)))
	mux.Handle("/api/photos", authMiddleware(http.HandlerFunc(listPhotosListHandler)))
	mux.Handle("/api/photos/upload", authMiddleware(http.HandlerFunc(uploadPhotoHandler)))
	mux.Handle("/api/photos/{id:[0-9]}", authMiddleware(http.HandlerFunc(getPhotoHandler)))
	mux.Handle("/api/photos/delete/", authMiddleware(http.HandlerFunc(deletePhotoHandler)))

	log.Println("Starting FamShare v0.1 on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
