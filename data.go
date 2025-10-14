package main

import "time"

// Data Structures
type User struct{ Username, Role string }
type Photo struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Owner    string    `json:"owner"`
	IsPublic bool      `json:"is_public"`
	Filename string    `json:"filename"`
	Date     time.Time `json:"date"`
}

// In-Memory Database
var users = map[string]User{
	"wissem": {Username: "wissem", Role: "admin"},
	"adam":   {Username: "adam", Role: "user"},
	"cal":    {Username: "cal", Role: "user"},
	"emy":    {Username: "emy", Role: "user"},
	"rita":   {Username: "rita", Role: "user"},
}

var photos = map[int]Photo{
	1: {ID: 1, Title: "Adam at the beach", Owner: "adam", IsPublic: true, Filename: "adam_beach.png", Date: time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC)},
	2: {ID: 2, Title: "Cal and Adam playing tennis", Owner: "cal", IsPublic: true, Filename: "cal_adam_tenis.png", Date: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)},
	3: {ID: 3, Title: "A day at the farm", Owner: "cal", IsPublic: true, Filename: "cal_farm.png", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
	4: {ID: 4, Title: "Emy with her newborn", Owner: "emy", IsPublic: false, Filename: "emy_baby.png", Date: time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC)},
	5: {ID: 5, Title: "Emy and Rita's Chess Game", Owner: "emy", IsPublic: true, Filename: "emy_rita_enjoying_chess.png", Date: time.Date(2024, 2, 18, 0, 0, 0, 0, time.UTC)},
	6: {ID: 6, Title: "Coffee with Emy", Owner: "wissem", IsPublic: false, Filename: "emy_uncle_coffee.png", Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)},
	7: {ID: 7, Title: "The Big Family Reunion", Owner: "wissem", IsPublic: true, Filename: "family_reunion.png", Date: time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC)},
	8: {ID: 8, Title: "Rita's Graduation", Owner: "emy", IsPublic: true, Filename: "rita_diploma.png", Date: time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC)},
	9: {ID: 9, Title: "Wissem's New Setup", Owner: "wissem", IsPublic: false, Filename: "uncle_new_setup.png", Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)},
}

var nextPhotoID = len(photos) + 1
