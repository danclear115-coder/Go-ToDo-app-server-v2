package handlers

import (
	"fmt"
	database "server/createDb"
	"server/functions"
	"server/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTask(c *fiber.Ctx) error {

	req := new(models.CreateTaskRequest)
	if err := c.BodyParser(&req); err != nil || req.Title == "" || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	userID, err := functions.AuthenticateUser(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect username and password",
		})
	}

	_, err = database.DB.Exec(
		"INSERT INTO tasks (title, content, priority, is_completed, user_id) VALUES (?, ?, ?, ?, ?)",
		req.Title, req.Content, req.Priority, false, userID,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON("Task has been created")

}

func GetTasks(c *fiber.Ctx) error {

	req := new(models.GetTasksRequest)

	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	userId, err := functions.AuthenticateUser(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect username and password",
		})
	}

	rows, err := database.DB.Query(
		"SELECT id, title, content, is_completed, priority, user_id FROM tasks WHERE user_id = ?",
		userId,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		var t models.Task

		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.IsCompleted, &t.Priority, &t.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		tasks = append(tasks, t)
	}

	return c.Status(fiber.StatusAccepted).JSON(tasks)

}

func CompleteTask(c *fiber.Ctx) error {

	req := new(models.CompleteTaskRequest)
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	userID, err := functions.AuthenticateUser(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect username and password",
		})
	}

	actualStatus := false

	err = database.DB.QueryRow(
		"SELECT is_completed FROM tasks WHERE user_id = ? AND id = ?",
		userID,
		req.Id,
	).Scan(&actualStatus)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "task not found or invalid data",
		})
	}

	_, err = database.DB.Exec(
		"UPDATE tasks SET is_completed = ? WHERE user_id = ? AND id = ?",
		!actualStatus,
		userID,
		req.Id,
	)

	return c.Status(fiber.StatusAccepted).JSON("Task has been completed")

}

func ChangeTask(c *fiber.Ctx) error {
	req := new(models.ChangeTaskRequest)
	if err := c.BodyParser(&req); err != nil || req.Content == "" || req.Username == "" || req.Password == "" {
		fmt.Println("Ошибка парсинга или пустые поля:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	userID, err := functions.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		fmt.Println("Ошибка авторизации пользователя:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect username and password",
		})
	}

	fmt.Printf("Отправка в БД -> Title: %s, Content: %s, Priority: %v, UserID: %v, TaskID: %v\n",
		req.Title, req.Content, req.Priority, userID, req.Id)

	res, err := database.DB.Exec(
		"UPDATE tasks SET title = ?, content = ?, priority = ? WHERE user_id = ? AND id = ?",
		req.Title,
		req.Content,
		req.Priority,
		userID,
		req.Id,
	)

	if err != nil {
		fmt.Println("Критическая ошибка SQL:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "server error",
		})
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println("Строк изменено в базе данных:", rowsAffected)

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "задача не найдена или данные идентичны старым",
		})
	}

	result := models.ModifiedTask{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	}

	return c.Status(fiber.StatusAccepted).JSON(result)
}

func DeleteTask(c *fiber.Ctx) error {

	req := new(models.DeleteTaskRequest)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid data or missing fields",
		})
	}

	userID, err := functions.AuthenticateUser(req.Username, req.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Incorrect username and password",
		})
	}

	var exists bool

	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM tasks WHERE user_id = ? AND id = ?)",
		userID,
		req.Id,
	).Scan(&exists)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "task not found"})
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Task dont find",
		})
	}

	_, err = database.DB.Exec(
		"DELETE from tasks WHERE id = ? AND user_id = ?",
		req.Id,
		userID,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "server error",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON("Task was been deleted")

}
