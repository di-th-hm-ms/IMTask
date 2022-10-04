package main

import (
	// "fmt"
	// "log"
	"time"
	// "net/http"
	// "os"
	// "database/sql"
	"github.com/gin-gonic/gin"
	// "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/cors"
	"IMTask/golang/src/controller"
)

func main() {

	engine := gin.Default()

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
			// taskEngine.POST("/add", controller.AddTask)
			// taskEngine.POST("/update", controller.UpdateTask)
			// taskEngine.POST("/delete", controller.DeleteTask)
		}
	}
	engine.Run(":8080")
}
