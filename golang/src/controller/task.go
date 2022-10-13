package controller

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"os"
	"IMTask/golang/src/service"
	"IMTask/golang/src/model"
	"github.com/gin-gonic/gin"
)

var DBSettingServiceIns service.DBSettingService = service.DBSettingService{}
var TaskServiceIns service.TaskService = service.TaskService{}
var UserServiceIns service.UserService = service.UserService{}

type JsonReq struct {
	Id			int64 `json:"field_id"`
	Title		string `json:"field_title"`
	IsAchieved 	bool `json:"field_isAchieved"`
	UserId 		string `json:"field_userId"`
}

const SUCCESS = 200
const BAD_REQUEST = 400

var DB *sql.DB
func InitDB() {
	ch := make(chan *sql.DB)
	go DBSettingServiceIns.OpenConnection(ch)
	DB = <- ch
	TaskServiceIns.DropTaskTable(DB)
	TaskServiceIns.CreateTaskTable(DB)
	UserServiceIns.DropUserTable(DB)
	UserServiceIns.CreateUserTable(DB)
}
func CloseDB() {
	DBSettingServiceIns.CloseConnection(DB)
}
func DropDB() {
	if _, err := DB.Exec("DROP DATABASE " + os.Getenv("MYSQL_DATABASE")); err != nil {
		log.Fatal(err)
	}
}
func GetTasks(c *gin.Context) {
	var tasks []model.Task = TaskServiceIns.GetTasksFromDB(DB)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": tasks,
	})
}

func AddTask(c *gin.Context) {
	var jsonReq JsonReq
	if err := c.ShouldBindJSON(&jsonReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}

	task := TaskServiceIns.InsertTaskIntoDB(jsonReq.Title, jsonReq.IsAchieved, jsonReq.UserId, DB)

	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		"Value": task,
	})
}

func UpdateTask(c *gin.Context) {
	var jsonReq JsonReq
	if err := c.ShouldBindJSON(&jsonReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}

	if err := TaskServiceIns.UpdateTaskOnDB(jsonReq.Id, jsonReq.Title, jsonReq.IsAchieved, jsonReq.UserId, DB); err != nil {
		// userId is wrong
		// c.AbortWithStatus(400)
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}
	task := model.Task{Id: jsonReq.Id, Title: jsonReq.Title, IsAchieved: jsonReq.IsAchieved, UserId: jsonReq.UserId}
	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		"Value": task,
	})
}


func DeleteTask(c *gin.Context) {
	var jsonReq JsonReq
	if err := c.ShouldBindJSON(&jsonReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task, err := TaskServiceIns.DeleteTaskFromDB(jsonReq.Id, jsonReq.UserId, DB); task == nil {
		// c.AbortWithStatus(204)
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Status": SUCCESS,
			"Value": task,
		})
	}
}
