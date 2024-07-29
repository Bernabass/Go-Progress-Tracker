package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
	"net/http"
)


type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
   }
   
var collection *mongo.Collection


func createTask(c *gin.Context) {
	var newTask Task
	err := c.BindJSON(&newTask)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	insertResult, err := collection.InsertOne(context.TODO(), newTask)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to create task"})
        return

	}

	c.IndentedJSON(http.StatusCreated, gin.H{
			"message": "task created successfully",
			"taskID":  insertResult.InsertedID,
		})
}

func getTasks(c *gin.Context) {
	var tasks []Task

	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch tasks"})

	}

	for cursor.Next(context.TODO()) {
		var task Task
		err := cursor.Decode(&task)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to decode task"})

		}

		tasks = append(tasks, task)
	}

    if err := cursor.Err(); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "cursor error"})
        return
    }

	cursor.Close(context.TODO())

    c.IndentedJSON(http.StatusOK, tasks)

}

func getTaskByID(c *gin.Context) {
	id := c.Param("id")

	var task Task
	
	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&task)

    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
        return
    }


	c.IndentedJSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask Task
	err := c.BindJSON(&updatedTask)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	
	}

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{
		{"title", updatedTask.Title},
		{"description", updatedTask.Description},
		{"status", updatedTask.Status}}}}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to update task"})
	}

	if updateResult.MatchedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
        return
    }

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated"})

}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	deleteResult, err := collection.DeleteOne(context.TODO(), bson.D{{"id", id}})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to delete task"})
		return
	}

	if deleteResult.DeletedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})

}

func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("TaskManagementDB").Collection("tasks")


	router := gin.Default()
	router.POST("/tasks", createTask)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.Run("localhost:8080")

}