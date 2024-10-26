package controllers

import (
	"be-recipe/src/config"
	"be-recipe/src/helpers"
	"be-recipe/src/models"
	"be-recipe/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AllRecipes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))

	return c.JSON(models.Paginate(config.DB, &models.Recipe{}, page))
}

func CreateRecipe(c *fiber.Ctx) error {
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

	uploadResult, err := services.UploadCloudinary(c, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	name := form.Value["Name"][0]
	ingredient := form.Value["Ingredient"][0]
	videourl := form.Value["VideoUrl"][0]

	userIDStr := form.Value["UserId"][0]
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return err
	}
	userIDUint := uint(userID)

	recipe := models.Recipe{
		Name:       name,
		Ingredient: ingredient,
		Photo:      uploadResult.SecureURL,
		VideoUrl:   videourl,
		UserId:     userIDUint,
	}

	config.DB.Create(&recipe)

	return c.JSON(fiber.Map{
		"Message": "Recipe created",
		"data":    recipe,
	})
}

func GetRecipe(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var recipe models.Recipe

	recipe.Id = uint(id)

	config.DB.Find(&recipe)

	return c.JSON(recipe)
}

func UpdateRecipe(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var recipe models.Recipe

	recipe.Id = uint(id)

	if err := c.BodyParser(&recipe); err != nil {
		return err
	}

	config.DB.Model(&recipe).Updates(recipe)

	return c.JSON(recipe)
}

func UpdatePhotoRecipe(c *fiber.Ctx) error {
	//update photo beserta isinya (foto harus terisi)
	id, _ := strconv.Atoi(c.Params("id"))

	var recipe models.Recipe

	if err := config.DB.First(&recipe, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Recipe not found")
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

	recipe.Photo = uploadResult.URL

	recipe.Id = uint(id)

	if err := c.BodyParser(&recipe); err != nil {
		return err
	}

	config.DB.Model(&recipe).Updates(recipe)

	return c.JSON(recipe)
}

func DeleteRecipe(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var recipe models.Recipe

	recipe.Id = uint(id)

	config.DB.Delete(&recipe)

	return c.JSON(fiber.Map{
		"Message": "Delete Complete",
	})
}
