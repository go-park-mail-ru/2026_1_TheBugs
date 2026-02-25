package validator

import (
	"regexp"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const maxEmailLength = 254

const minPwdLenght = 8
const maxPwdLenght = 64

func ValidateEmail(email string) bool {
	if len(email) > maxEmailLength {
		return false
	}
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

func ValidatePwd(pwd string) bool {
	if len(pwd) > maxPwdLenght || len(pwd) < minPwdLenght {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pwd)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pwd)
	hasDigit := regexp.MustCompile(`\d`).MatchString(pwd)
	return hasUpper && hasLower && hasDigit
}
