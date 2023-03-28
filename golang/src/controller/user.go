package controller

import (
	"net/http"
	// "database/sql"
	// "os"
	"IMTask/golang/src/service"
	"IMTask/golang/src/handler"
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
		"Status": model.SUCCESS,
		"Value": users,
	})
}

func GetUser(c *gin.Context) {
	type UserAuth struct {
		// Id string
		Email 	 string
		Password string
	}
	var userAuth UserAuth
	if err := c.ShouldBindJSON(&userAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "Status": BAD_REQUEST,
			"Value": err.Error(),
		})
	}
	if user, serverError := UserServiceIns.GetUserFromDB(userAuth.Email, DB); serverError != nil {
		c.JSON(serverError.Status, gin.H{
			// "Status": serverError.Status,
			"Value": serverError.Msg,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			// "Status": SUCCESS,
			"data": user,
		})
	}
}
// DEBUG

func AddUser(c *gin.Context) {
	// var userReq model.UserReq
	content := model.NewUserContent()
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "Status": BAD_REQUEST,
			"value": err.Error(),
		})
		return
	}

	// user := UserServiceIns.InsertUserIntoDB(userReq.Email, userReq.Username, userReq.Password, userReq.CreatedAt, DB)
	user, serverErr := UserServiceIns.InsertUserIntoDB(content, DB)
	if serverErr != nil {
		c.JSON(serverErr.Status, gin.H{
			// "Status": err.Status,
			"value": serverErr.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": model.SUCCESS,
		"value": user,
	})


}

func UpdateUsername(c *gin.Context) {
	userReq := model.NewUserReq()
	if err := c.ShouldBindJSON(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "status": BAD_REQUEST,
			"value": err.Error(),
		})
		return
	}

	if err := UserServiceIns.UpdateUsernameOnDB(userReq, DB); err != nil {
		c.JSON(err.Status, gin.H{
			// "status": err.Status,
			"value": err.Msg,
		})
		return
	}
	user := model.NewUser()
	user.Bind(userReq)
	c.JSON(http.StatusOK, gin.H{
		"status": model.SUCCESS,
		"value": user,
	})
}


func DeleteUser(c *gin.Context) {
	userReq := model.NewUserReq()
	if err := c.ShouldBindJSON(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "Status": BAD_REQUEST,
			"value":  err.Error(),
		})
		return
	}

	if err := UserServiceIns.DeleteUserFromDB(userReq, DB); err != nil {
		c.JSON(err.Status, gin.H{
			"value": err.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": model.SUCCESS,
		"value": "deleted",
	})
}

func Login(c *gin.Context) {
	loginSt := &struct {
		Email		string `json:"email"`
		// Username	string `json:"username"`
		Password	string `json:"password"`
	}{}
	if err := c.ShouldBindJSON(loginSt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "Status": BAD_REQUEST,
			"value":  err.Error(),
		})
		return
	}
	user, err := UserServiceIns.Login(loginSt.Email, loginSt.Password, DB);
	if err != nil {
		c.JSON(err.Status, gin.H{
			// "Status": err.Status,
			"value": err.Msg,
		})
		return
	}
	token := handler.GenerateJWT(user.Id, user.Username)
	c.JSON(http.StatusOK, gin.H{
		"status": model.SUCCESS,
		"value": token,
	})
}

// func Signup(c *gin.Context) {
// 	content := model.NewUserContent()
// 	if err := c.ShouldBindJSON(content); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"Status": BAD_REQUEST,
// 			"Value":  err.Error(),
// 		})
// 		return
// 	}
// 	user, err := UserServiceIns.InsertUserIntoDB(content, DB)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"Status": err.Status,
// 			"Value": err.Msg,
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"Status": SUCCESS,
// 		"Value": user,
// 	})
// }
