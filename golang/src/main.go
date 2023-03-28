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
	pageHub := handler.NewHub()
	go pageHub.Loop()
	go func() {
		http.HandleFunc("/ws", handler.New(pageHub).Handle)
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
			// handler.JwtMiddleWare.
			taskEngine.GET("/list", handler.VerifyJWT(&gin.Context, controller.GetTasks))

			taskEngine.POST("/add", controller.AddTask)
			taskEngine.POST("/update", controller.UpdateTask)
			taskEngine.POST("/delete", controller.DeleteTask)
		}
		userEngine := v1.Group("/users")
		{
			userEngine.POST("/login", controller.Login)
			userEngine.POST("/signup", controller.AddUser)
			userEngine.POST("/", controller.GetUser)
			// TODO specific user
			userEngine.GET("/list", controller.GetUsers)
			// userEngine.POST("/add", controller.AddUser)
			// middleware check
			userEngine.POST("/update-un", controller.UpdateUsername)
			userEngine.POST("/delete", controller.DeleteUser)
		}
	}
	engine.Run(":8080")
}
