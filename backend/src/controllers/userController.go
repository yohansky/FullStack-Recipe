package controllers

import (
	"be-recipe/src/config"
	"be-recipe/src/helpers"
	"be-recipe/src/models"
	"be-recipe/src/services"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(config.DB, &models.User{}, page))
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var user models.User

	user.Id = uint(id)

	if err := config.DB.Find(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ID tidak ditemukan",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.JSON(fiber.Map{
			"Message": "Id tidak ditemukan",
		})
	}

	var user models.User

	user.Id = uint(id)

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	config.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func UpdatePhotoUser(c *fiber.Ctx) error {
	//update photo beserta isinya (foto harus terisi)
	id, _ := strconv.Atoi(c.Params("id"))

	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	file, err := c.FormFile("Photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Gagal mengunggah file: " + err.Error())
	}

	maxFileSize := int64(2 << 20)
	if err := helpers.SizeUploadValidation(file.Size, maxFileSize); err != nil {
		return err
	}

	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal membuka file: " + err.Error())
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	if _, err := fileHeader.Read(buffer); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Gagal membaca file: " + err.Error())
	}

	validFileTypes := []string{"image/png", "image/jpeg", "image/jpg", "application/pdf"}
	if err := helpers.TypeUploadValidation(buffer, validFileTypes); err != nil {
		return err
	}

	fileHeader.Seek(0, 0)

	uploadResult, err := services.UploadCloudinary(c, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user.Photo = uploadResult.URL

	user.Id = uint(id)

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	config.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var user models.User

	user.Id = uint(id)

	config.DB.Delete(&user)

	return c.JSON(fiber.Map{
		"Message": "Delete Complete",
	})
}
