package controller

import (
	"fmt"
	// "log"
	"net/http"
	// "database/sql"
	// "os"
	"IMTask/golang/src/service"
	"IMTask/golang/src/model"
	"github.com/gin-gonic/gin"
)




var UserService service.UserService = service.UserService{}

func GetUsers(c *gin.Context) {
	var users []model.User = UserServiceIns.GetUsersFromDB(DB)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": users,
	})
}
func GetUser(c *gin.Context) {
	type UserAuth struct {
		Id string
		Password string
	}
	var userAuth UserAuth
	if err := c.ShouldBindJSON(&userAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
	}
	var user = UserServiceIns.GetUserFromDB(userAuth.Id, userAuth.Password, DB)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": user,
	})
}

func AddUser(c *gin.Context) {
	var userReq model.UserReq
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}

	// user := UserServiceIns.InsertUserIntoDB(userReq.Email, userReq.Username, userReq.Password, userReq.CreatedAt, DB)
	if user, err := UserServiceIns.InsertUserIntoDB(&userReq, DB); user != nil {
		c.JSON(http.StatusOK, gin.H{
			"Status": SUCCESS,
			"Value": user,
		})
	} else if err.Error() == "400" {
		// email error

		// TODO password error
		// req error
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": "parameter error",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
	}


}

// TODO JWT開封後一致後
func UpdateUser(c *gin.Context) {
	userReq := model.NewUserReq()
	if err := c.ShouldBindJSON(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}

	if err := UserServiceIns.UpdateUserOnDB(userReq, DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}
	user := model.NewUser()
	user.Bind(userReq)
	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		"Value": user,
	})
}


func DeleteUser(c *gin.Context) {
	userReq := model.NewUserReq()
	if err := c.ShouldBindJSON(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UserServiceIns.DeleteUserFromDB(userReq, DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		// "Value": user,
	})
}
