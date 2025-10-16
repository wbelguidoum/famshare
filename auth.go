package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const (
	SESSION_COOKIE   string     = "user_cookie"
	USER_CONTEXT_KEY contextKey = "user"
)

func loginHandler(c *fiber.Ctx) error {
	var creds struct {
		Username string `json:"username"`
	}
	if err := c.BodyParser(&creds); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request format")
	}

	if _, ok := userStore.Get(creds.Username); !ok {
		return fiber.ErrUnauthorized
	}

	c.Cookie(&fiber.Cookie{
		Name:     SESSION_COOKIE,
		Value:    creds.Username,
		Expires:  time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func logoutHandler(c *fiber.Ctx) error {
	// Expire the cookie by setting its expiration date to the past.
	c.Cookie(&fiber.Cookie{
		Name:     SESSION_COOKIE,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to the past
		HTTPOnly: true,
		Path:     "/",
	})
	return c.SendStatus(fiber.StatusOK)
}

func checkSessionHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(user)
}

// authMiddleware my insecure middleware :)
func authMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies(SESSION_COOKIE)
	if cookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: No cookie found")
	}

	// We trust the username value directly from the cookie.
	username := cookie
	if _, ok := userStore.Get(username); !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Invalid user")
	}

	c.Locals(USER_CONTEXT_KEY, username)

	return c.Next()
}

func getConnectedUser(c *fiber.Ctx) (*User, error) {
	username, ok := c.Locals(USER_CONTEXT_KEY).(string)
	if !ok {
		return nil, fmt.Errorf("could not retrieve user information from context")
	}

	user, userExists := userStore.Get(username)
	if !userExists {
		return nil, fmt.Errorf("user specified in context not found")
	}

	return &user, nil
}
