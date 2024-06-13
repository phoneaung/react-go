package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// bson is binary json format that mongodb uses
type Todo struct {
	ID        int    `json:"_id" bson:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!", err)
	}

	// get mongodb uri from .env
	MONGODB_URI := os.Getenv("MONGODB_URI")

	// connect to mongo db connection string
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	// get client and error out of this call
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos/:id", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos:id", deleteTodos)

	// listen to port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen("0.0.0.0" + port))

}
