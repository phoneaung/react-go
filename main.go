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
	// application # boilerplate
	app := fiber.New()

	todos := []Todo{}

	// handler function # boilerplate
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// create a todo
	app.Post("/api/todos/create", func(c *fiber.Ctx) error {
		todo := &Todo{} // create a todo with default Todo struct values id=0, completed=false, Body=""

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		// if the user does not type anything in body, return an error
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"msg": "error: todo body is required!"})
		}

		// increment id numbers
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// update a todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(400).JSON(fiber.Map{"error": "Todo Not found"})
	})

	// delete a todo

	// listening to port # boilerplate
	log.Fatal(app.Listen(":4000"))
}
