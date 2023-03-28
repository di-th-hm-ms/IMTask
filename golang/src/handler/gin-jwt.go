package handler

import (
	// "crypto/rsa"
	"log"
	"os"
	"errors"
	"net/http"
	"time"

	// "IMTask/golang/src/model"
	"github.com/gin-gonic/gin"
	gojwt "github.com/form3tech-oss/jwt-go"
	"github.com/golang-jwt/jwt"
)

// jwt
// var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
func GenerateJWT(userId, username string) string {
	// header
	// token := gojwt.New(gojwt.SigningMethodHS256)
	token := gojwt.New(gojwt.SigningMethodHS256)

	claims := token.Claims.(gojwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = userId
	claims["name"] = username
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenStr, _ := token.SignedString([]byte(os.Getenv("SIGNIN_KEY")))

	// w.Write([]byte(tokenStr))
	return tokenStr
}

// func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
// func verifyJWT(c *gin.Context) http.HandlerFunc {
func VerifyJWT(c *gin.Context, handler func(c *gin.Context)) func(c *gin.Context) {
	log.Println(c.Request.Header)
	if headerToken := c.Request.Header["Token"]; headerToken != nil {
		// bearer
		log.Println("header token")
		log.Println(headerToken)
		token, err := jwt.Parse(headerToken[0], func(token *jwt.Token) (interface{}, error) {
			// if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"value": "You're not unauthorized",
				})
				return "", errors.New("Auth error")
			}
			return "", nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"value": "You're not unauthorized",
			})
			return func(c *gin.Context) {
				log.Println("unauthorized")
			}
		} else {
			return handler
		}
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"value": "BAD REQUEST",
		})
		log.Println("Not found Header token")
	}
}
