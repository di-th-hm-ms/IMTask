package handler

import (
	// "os"
	// "time"
	"golang.org/x/crypto/bcrypt"
)

// bcrypt
func Encrypto(password string) string {
	// encrypto
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password),10)
	return string(hashed)
}

func Compare(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}




// var JwtMiddleWare = jwtmiddleware.New(jwtmiddleware.Options{
// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("SIGNIN_KEY")), nil
// 	},
// 	SigningMethod: jwt.SigningMethodHS256,
// })
