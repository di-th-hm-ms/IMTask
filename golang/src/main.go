package main

import (
	// "fmt"
	"log"
	"time"
	"net/http"
	// "os"
	// "database/sql"
	"github.com/gin-gonic/gin"
	// "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/cors"
	"IMTask/golang/src/controller"
	"IMTask/golang/src/handler"
)

func main() {
	// DB
	controller.InitDB()
	defer controller.CloseDB()

	//ws handle
	// mux := http.
	go func() {
		http.HandleFunc("/ws", handler.New().Handle)
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Printf("ws error: %s", err)
		}
	}()

	engine := gin.Default()
	// middleware
	// engine.Use()

	// cors before routing
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:4200",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			// "OPTIONS"
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge: 24 * time.Hour,
	}))

	// routing
	v1 := engine.Group("/v1")
	{
		taskEngine := v1.Group("/tasks")
		{
			taskEngine.GET("/list", controller.GetTasks)
			taskEngine.POST("/add", controller.AddTask)
			taskEngine.POST("/update", controller.UpdateTask)
			taskEngine.POST("/delete", controller.DeleteTask)
		}
		userEngine := v1.Group("/users")
		{
			userEngine.POST("/", controller.GetUser) // TODO specific user
			userEngine.GET("/list", controller.GetUsers)
			userEngine.POST("/add", controller.AddUser)
			userEngine.POST("/update", controller.UpdateUser)
			userEngine.POST("/delete", controller.DeleteUser)
		}
	}
	engine.Run(":8080")
}
