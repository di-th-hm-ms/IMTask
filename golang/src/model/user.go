package model

import (
	// "github.com/gin-gonic/gin"
	"crypto/rand"
)
type User struct {
	Id			string
	Email 		string
	Username	string
	Password 	string
	CreatedAt	string
}
type UserReq struct {
	Id			string `json:"f_user_id"`
	Email 		string `json:"f_email"`
	Username	string `json:"f_username"`
	Password 	string `json:"f_password"`
	CreatedAt	string `json:"f_created_at"`
}
type UserContent struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
	Password	string `json:"password"`
}

func (u *User) Bind(userReq *UserReq) {
	u.Id = userReq.Id
	u.Email = userReq.Email
	u.Username = userReq.Username
	u.Password = userReq.Password
	u.CreatedAt = userReq.CreatedAt
}

func NewUser() *User {
	return &User{}
}
func NewUserReq() *UserReq {
	return &UserReq{}
}
func NewUserContent() *UserContent{
	return &UserContent{}
}

func NewFromUserReq(userReq *UserReq) *User {
	user := User {
		Id: userReq.Id,
		Email: userReq.Email,
		Username: userReq.Username,
		Password: userReq.Password,
		CreatedAt: userReq.CreatedAt,
	}
	return &user
}

func GenerateRandStr(size uint32) (result string, err error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bs := make([]byte, size)
	result = ""
	if _, err = rand.Read(bs); err != nil {
		return
	}
	for _, b := range bs {
		result += string(letters[int(b)%len(letters)])
	}
	return
}
