package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// bson is binary json format that mongodb uses
// mongoDB has its own datatype so instead of int for ID, we use primitive.ObjectID which is from mongo
// by default, ID is gonna be zero for first todo created, so we have to omit the value
type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
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

	// when main is done executing, we will disconnect from the database
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas!")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos/:id", createTodo)
	// app.Patch("/api/todos/:id", updateTodo)
	// app.Delete("/api/todos:id", deleteTodo)

	// listen to port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// get all todos
func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	// no filter, get them all todos from mongodb
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	// defer is cure to use postpone function until the surroundings end or completed
	defer cursor.Close(context.Background())

	// iterate through the cursor to decode each todo item
	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err // return an error if decoding fails
		}
		// append each decoded todo
		todos = append(todos, todo)
	}
	// return the todos slice as json format
	return c.JSON(todos)
}

// create a todo
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	// if a new todo has empty string, return error
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Todo cannot be an empty!"})
	}

	// save todo to database (created a todo)
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	// Object ID
	todo.ID := insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}
