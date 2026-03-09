package validator

import (
	"regexp"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const pwdRegexPattern = `^[a-zA-Z\d!@#$%^&*\-]{8,}$`
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
	neededSymbols := regexp.MustCompile(pwdRegexPattern).MatchString(pwd)
	return hasUpper && hasLower && hasDigit && neededSymbols
}

func ValidateCred(email string, pwd string) error {
	if !ValidateEmail(email) {
		return entity.NewValidationError("email")
	}
	if !ValidatePwd(pwd) {
		return entity.NewValidationError("password")
	}
	return nil
}
