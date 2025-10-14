package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

const PHOTO_DIR = "gallery"

func listPhotosHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		log.Println("[listPhotosHandler] Error retrieving user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	visiblePhotos := []Photo{}
	allPhotos := photoStore.GetAll()

	for _, photo := range allPhotos {
		if photo.IsPublic || photo.Owner == user.Username {
			visiblePhotos = append(visiblePhotos, photo)
		}
	}

	return c.JSON(visiblePhotos)
}

func getPhotoHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		log.Println("[getPhotoHandler] Error retrieving user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("[getPhotoHandler] Error parsing photo ID:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid photo ID")
	}

	photo, found := photoStore.Get(id)
	if !found {
		log.Println("[getPhotoHandler] Photo not found for ID:", id)
		return fiber.ErrNotFound
	}

	if !photo.IsPublic && photo.Owner != user.Username {
		log.Printf("[getPhotoHandler] Access denied: photo is public %v, owner %v, user %v\n", photo.IsPublic, photo.Owner, user.Username)
		return fiber.ErrForbidden
	}

	filePath := filepath.Join(PHOTO_DIR, photo.Filename)
	return c.SendFile(filePath)
}

// updatePhotoHandler handles full updates to a photo resource.
func updatePhotoHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		log.Println("[updatePhotoHandler] Error retrieving user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("[updatePhotoHandler] Error parsing photo ID:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid photo ID")
	}

	originalPhoto, found := photoStore.Get(id)
	if !found {
		log.Println("[updatePhotoHandler] Photo not found for ID:", id)
		return fiber.ErrNotFound
	}

	if originalPhoto.Owner != user.Username {
		log.Printf("[updatePhotoHandler] Access denied: photo owner %v, user %v\n", originalPhoto.Owner, user.Username)
		return fiber.ErrForbidden
	}

	var updatedPhoto Photo
	if err := c.BodyParser(&updatedPhoto); err != nil {
		log.Println("[updatePhotoHandler] Error parsing request body:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	updatedPhoto.ID = originalPhoto.ID
	updatedPhoto.Owner = originalPhoto.Owner
	updatedPhoto.Filename = originalPhoto.Filename
	updatedPhoto.Date = originalPhoto.Date

	savedPhoto, ok := photoStore.Update(updatedPhoto)
	if !ok {
		log.Println("[updatePhotoHandler] Failed to update photo for ID:", id)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update photo")
	}

	return c.JSON(savedPhoto)
}

func uploadPhotoHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		log.Println("[uploadPhotoHandler] Error retrieving user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	title := c.FormValue("title")
	if title == "" {
		log.Println("[uploadPhotoHandler] Title is required")
		return fiber.NewError(fiber.StatusBadRequest, "Title is required")
	}

	isPrivateValue := c.FormValue("is_private")

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		log.Println("[uploadPhotoHandler] Error retrieving the file:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error retrieving the file")
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	if err := c.SaveFile(fileHeader, filepath.Join(PHOTO_DIR, filename)); err != nil {
		log.Println("[uploadPhotoHandler] Error saving the file:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error saving the file")
	}

	newPhoto := Photo{
		Title:    title,
		Owner:    user.Username,
		IsPublic: isPrivateValue != "on",
		Filename: filename,
		Date:     time.Now(),
	}
	id := photoStore.Add(newPhoto)
	newPhoto.ID = id

	return c.Status(fiber.StatusCreated).JSON(newPhoto)
}

func deletePhotoHandler(c *fiber.Ctx) error {
	user, err := getConnectedUser(c)
	if err != nil {
		log.Println("[deletePhotoHandler] Error retrieving user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		log.Println("[deletePhotoHandler] Error parsing photo ID:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid photo ID")
	}

	photo, found := photoStore.Get(id)
	if !found {
		log.Println("[deletePhotoHandler] Photo not found for ID:", id)
		return fiber.ErrNotFound
	}

	// Security check: Only the owner or an admin can delete.
	if photo.Owner != user.Username && user.Role != "admin" {
		log.Printf("[deletePhotoHandler] Access denied: photo owner %v, user %v (admin? %v)\n", photo.Owner, user.Username, user.Role == "admin")
		return fiber.ErrForbidden
	}

	photoStore.Delete(id)

	return c.JSON(fiber.Map{"message": "Photo deleted successfully"})
}
