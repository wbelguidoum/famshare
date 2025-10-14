package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func listPhotosListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

func getPhotoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	photos
}

func uploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	newPhoto := Photo{
		ID:       nextPhotoID,
		Title:    title,
		Owner:    "developer",
		IsPublic: true,
		Filename: filename,
	}
	photos = append(photos, newPhoto)
	nextPhotoID++
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPhoto)
}

func deletePhotoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := filepath.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}
	var photoIndex = -1
	var photoFilename string
	for i, p := range photos {
		if p.ID == id {
			photoIndex = i
			photoFilename = p.Filename
			break
		}
	}
	if photoIndex == -1 {
		http.NotFound(w, r)
		return
	}
	os.Remove(filepath.Join("uploads", photoFilename))
	photos = append(photos[:photoIndex], photos[photoIndex+1:]...)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Photo deleted successfully"})
}
