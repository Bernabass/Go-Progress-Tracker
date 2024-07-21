package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
   }
   

   var tasks = []Task{
	   {ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	   {ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	   {ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
   }

func createTask(c *gin.Context) {
	var newTask Task
	err := c.BindJSON(&newTask)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

func getTaskByID(c *gin.Context) {
	id := c.Param("id")

	for _, T := range tasks{
		if T.ID == id{
			c.IndentedJSON(http.StatusOK, T)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask Task
	err := c.BindJSON(&updatedTask)

	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	
	}

	for i, T := range tasks{
		if T.ID == id{
			if updatedTask.Title != ""{
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != ""{
				tasks[i].Description = updatedTask.Description
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "task updated"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	for i, T := range tasks {
		if T.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

func main() {
	router := gin.Default()
	router.POST("/tasks", createTask)
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.Run("localhost:8080")

}