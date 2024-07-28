package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` 
}

var jwtKey = []byte("shalamagando")

var taskCollection *mongo.Collection
var userCollection *mongo.Collection

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createUser(c *gin.Context) {
	var newUser User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var existingUser User
	err = userCollection.FindOne(context.TODO(), bson.D{{"username", newUser.Username}}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username already taken"})
		return
	}

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}

	newUser.Password = hashedPassword

	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func login(c *gin.Context) {
	var loginUser User
	err := c.BindJSON(&loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var user User
	err = userCollection.FindOne(context.TODO(), bson.D{{"username", loginUser.Username}}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	if !CheckPasswordHash(loginUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	tokenString := authHeader[len("Bearer "):]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
	c.Next()
	}


func adminMiddleware(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}
		c.Next()
	}


func createTask(c *gin.Context) {
	var newTask Task
	err := c.BindJSON(&newTask)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	insertResult, err := taskCollection.InsertOne(context.TODO(), newTask)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created successfully",
		"taskID":  insertResult.InsertedID,
	})
}

func getTasks(c *gin.Context) {
	var tasks []Task

	cursor, err := taskCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch tasks"})
		return
	}

	for cursor.Next(context.TODO()) {
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to decode task"})
			return
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cursor error"})
		return
	}

	cursor.Close(context.TODO())
	c.JSON(http.StatusOK, tasks)
}

func getTaskByID(c *gin.Context) {
	id := c.Param("id")

	var task Task
	err := taskCollection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask Task
	err := c.BindJSON(&updatedTask)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{
		{"title", updatedTask.Title},
		{"description", updatedTask.Description},
		{"status", updatedTask.Status},
	}}}

	updateResult, err := taskCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update task"})
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	deleteResult, err := taskCollection.DeleteOne(context.TODO(), bson.D{{"id", id}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete task"})
		return
	}

	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
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

	taskCollection = client.Database("TaskManagementDB").Collection("tasks")
	userCollection = client.Database("TaskManagementDB").Collection("users")

	router := gin.Default()
	router.POST("/register", createUser)
	router.POST("/login", login)

	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByID)

	router.PUT("/tasks/:id", authMiddleware, adminMiddleware, updateTask)
	router.POST("/tasks", authMiddleware, adminMiddleware, createTask)
	router.DELETE("/tasks/:id", authMiddleware, adminMiddleware, deleteTask)

	router.Run("localhost:8080")
}