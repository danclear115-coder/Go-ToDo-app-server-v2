package main

import (
	"log"
	database "server/createDb"
	"server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.InitDB()
	defer database.DB.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/user/createUser", handlers.CreateUser)
	app.Post("/user/login", handlers.UserLogin)

	app.Post("/tasks/createTask", handlers.CreateTask)
	app.Put("/tasks/completeTask", handlers.CompleteTask)
	app.Put("/tasks/changeTask", handlers.ChangeTask)
	app.Post("/tasks/getTasks", handlers.GetTasks)
	app.Post("/tasks/deleteTask", handlers.DeleteTask)

	log.Fatal(app.Listen(":3000"))
}
