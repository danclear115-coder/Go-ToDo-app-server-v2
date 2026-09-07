package handlers

import (
    "fmt"
	database "server/createDb"
	"server/models"
    "server/functions"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {

	req := new(models.CreateUserRequest)
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	res, err := database.DB.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		req.Username, string(hash),
	)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user already exists",
		})
	}

	id, _ := res.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(models.User{
		ID:       int(id),
		Username: req.Username,
	})

}

func UserLogin(c *fiber.Ctx) error {

	req := new(models.LoginRequest)

    if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data",
		})
	}

    if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username and password are required",
		})
	}

    userID, err := functions.AuthenticateUser(
		req.Username,
		req.Password,
	)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid username or password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user_id": userID,
	})

}
