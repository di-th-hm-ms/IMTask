package model
import (
	"regexp"
	"unicode"
)
func ValidateEmail(email string) bool {
	// r, _ := regexp.Compile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:.[a-zA-Z0-9-]+)*$")
	if len(email) > 50 { return false }
	r, _ := regexp.Compile(`[\w\-._]+@[\w\-._]+\.[A-Za-z]+`)
	return r.MatchString(email)
}

func ValidatePassword(s string) bool {
	// letters := 0
	passedCnt := 0
	number, upper, lower := false, false, false
    for _, c := range s {
        switch {
        case unicode.IsNumber(c):
			number = true
			passedCnt++
        case unicode.IsUpper(c):
			upper = true
			passedCnt++
        case unicode.IsLetter(c):
			lower = true
			passedCnt++
		default:
			return false
        }
	}
	return 8 <= passedCnt &&
		   passedCnt <= 100 &&
		   number &&
		   upper &&
		   lower &&
		   passedCnt == len(s)
}
