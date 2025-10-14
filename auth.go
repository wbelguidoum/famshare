package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type contextKey string

const (
	SESSION_COOKIE   string     = "user_cookie"
	USER_CONTEXT_KEY contextKey = "user"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user exists
	_, ok := users[creds.Username]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Set the cookie directly to the username.
	// This is INSECURE because a user can easily change their cookie value.
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE,
		Value:    creds.Username,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// "Logging out" now just means telling the browser to delete the cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE,
		Value:    "",
		Expires:  time.Now(), // Set expiration to the past
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
}

// checkSessionHandler checks if the user is logged in
func checkSessionHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER_CONTEXT_KEY).(*User)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// authMiddleware my insecure middleware :)
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SESSION_COOKIE)
		if err != nil {
			http.Error(w, "Unauthorized: No cookie found", http.StatusUnauthorized)
			return
		}

		// We trust the username value directly from the cookie.
		username := cookie.Value
		user, ok := users[username]
		if !ok {
			http.Error(w, "Unauthorized: Invalid user", http.StatusUnauthorized)
			return
		}

		// The user is "valid" because their name is in our list.
		// Add the username to the context for handlers to use.
		ctx := context.WithValue(r.Context(), USER_CONTEXT_KEY, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
