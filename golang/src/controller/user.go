package controller

import (
	"net/http"
	// "database/sql"
	// "os"
	"IMTask/golang/src/service"
	"IMTask/golang/src/model"
	"github.com/gin-gonic/gin"
)




var UserService service.UserService = service.UserService{}

// DEBUG
func GetUsers(c *gin.Context) {
	users, serverErr := UserServiceIns.GetUsersFromDB(DB)
	if serverErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"Status": serverErr.Status,
			"Value": serverErr.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		"Value": users,
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
	if user, serverError := UserServiceIns.GetUserFromDB(userAuth.Id, userAuth.Password, DB); serverError != nil {
		c.JSON(http.StatusOK, gin.H{
			"Status": serverError.Status,
			"Value": serverError.Msg,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Status": SUCCESS,
			"data": user,
		})
	}
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": err.Status,
			"Value": err.Msg,
		})
	}


}

func UpdateUser(c *gin.Context) {
	userReq := model.NewUserReq()
	if err := c.ShouldBindJSON(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": BAD_REQUEST,
			"Value": err.Error(),
		})
		return
	}

	if err := UserServiceIns.UpdateUsernameOnDB(userReq, DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": err.Status,
			"Value": err.Msg,
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
			"Status": err.Status,
			"Value": err.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": SUCCESS,
		// "Value": user,
	})
}
