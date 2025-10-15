package main

import (
	"sort"
	"sync"
	"time"
)

// Types
type User struct{ Username, Role string }
type Photo struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Owner    string    `json:"owner"`
	IsPublic bool      `json:"is_public"`
	Filename string    `json:"filename"`
	Date     time.Time `json:"date"`
}
type UserStore interface {
	Get(username string) (User, bool)
}
type PhotoStore interface {
	Get(id int) (Photo, bool)
	GetAll() []Photo
	Add(photo Photo) int
	Update(updated Photo) (Photo, bool)
	Delete(id int)
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
	1:  {ID: 1, Title: "Adam at the beach", Owner: "adam", IsPublic: true, Filename: "adam_beach.png", Date: time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC)},
	2:  {ID: 2, Title: "Cal and Adam playing tennis", Owner: "cal", IsPublic: true, Filename: "cal_adam_tenis.png", Date: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)},
	3:  {ID: 3, Title: "A day at the farm", Owner: "cal", IsPublic: true, Filename: "cal_farm.png", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
	4:  {ID: 4, Title: "Emy with her newborn", Owner: "emy", IsPublic: false, Filename: "emy_baby.png", Date: time.Date(2024, 7, 23, 0, 0, 0, 0, time.UTC)},
	5:  {ID: 5, Title: "Emy and Rita's Chess Game", Owner: "emy", IsPublic: true, Filename: "emy_rita_enjoying_chess.png", Date: time.Date(2024, 2, 18, 0, 0, 0, 0, time.UTC)},
	6:  {ID: 6, Title: "Coffee with Emy", Owner: "wissem", IsPublic: false, Filename: "emy_wissem_coffee.png", Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)},
	7:  {ID: 7, Title: "The Big Family Reunion", Owner: "wissem", IsPublic: true, Filename: "family_reunion.png", Date: time.Date(2024, 7, 22, 0, 0, 0, 0, time.UTC)},
	8:  {ID: 8, Title: "Rita's Graduation", Owner: "rita", IsPublic: true, Filename: "rita_diploma.png", Date: time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC)},
	9:  {ID: 9, Title: "Wissem's New Setup", Owner: "wissem", IsPublic: false, Filename: "wissem_new_setup.png", Date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)},
	10: {ID: 10, Title: "Wissem's Profile photo", Owner: "wissem", IsPublic: false, Filename: "wissem.png", Date: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)},
	11: {ID: 11, Title: "Emy's Profile photo", Owner: "emy", IsPublic: false, Filename: "emy.png", Date: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)},
	12: {ID: 12, Title: "Cal's Profile photo", Owner: "cal", IsPublic: false, Filename: "cal.png", Date: time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC)},
	13: {ID: 13, Title: "Adam's Profile photo", Owner: "adam", IsPublic: false, Filename: "adam.png", Date: time.Date(2021, 11, 2, 0, 0, 0, 0, time.UTC)},
	14: {ID: 14, Title: "Rita's Profile photo", Owner: "rita", IsPublic: false, Filename: "rita.png", Date: time.Date(2021, 11, 2, 0, 0, 0, 0, time.UTC)},
}

type InMemoryUserStore struct {
	Users map[string]User
	lock  sync.RWMutex
}

func (us *InMemoryUserStore) Get(username string) (User, bool) {
	us.lock.RLock()
	defer us.lock.RUnlock()
	user, ok := us.Users[username]
	return user, ok
}

type InMemoryPhotoStore struct {
	Photos      map[int]Photo
	nextPhotoID int
	lock        sync.RWMutex
}

func (ps *InMemoryPhotoStore) Get(id int) (Photo, bool) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	photo, ok := ps.Photos[id]
	return photo, ok
}

func (ps *InMemoryPhotoStore) GetAll() []Photo {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	photos := make([]Photo, 0, len(ps.Photos))
	for _, photo := range ps.Photos {
		photos = append(photos, photo)
	}
	sort.Slice(photos, func(i int, j int) bool {
		return photos[i].Date.After(photos[j].Date)
	})
	return photos
}

func (ps *InMemoryPhotoStore) Add(photo Photo) int {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	photo.ID = ps.nextPhotoID
	ps.Photos[photo.ID] = photo
	ps.nextPhotoID++
	return photo.ID
}

func (ps *InMemoryPhotoStore) Update(updated Photo) (Photo, bool) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	old, found := ps.Photos[updated.ID]
	if !found {
		return Photo{}, false
	}
	updated.Filename = old.Filename
	updated.Date = old.Date
	updated.Owner = old.Owner
	ps.Photos[updated.ID] = updated
	return updated, true
}

func (ps *InMemoryPhotoStore) Delete(id int) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	delete(ps.Photos, id)
}

var userStore UserStore = &InMemoryUserStore{Users: users}
var photoStore PhotoStore = &InMemoryPhotoStore{Photos: photos, nextPhotoID: len(photos) + 1}
