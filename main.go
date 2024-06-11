package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// custom kinda like datastructure

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World!")
	// application # boilerplate
	app := fiber.New()

	todos := []Todo{}

	// handler function # boilerplate
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "HelloWorld!"})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // create a todo with default Todo struct values id=0, completed=false, Body=""

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "error: todo body is required!"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// listening to port # boilerplate
	log.Fatal(app.Listen(":4000"))
}
