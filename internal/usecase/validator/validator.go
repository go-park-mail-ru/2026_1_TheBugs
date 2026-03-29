package validator

import (
	"mime/multipart"
	"path/filepath"
	"regexp"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const pwdRegexPattern = `^[a-zA-Z\d!@#$%^&*\-]{8,}$`

const phoneRegexp = `^(\+7|8)\s?[\s(]?\d{3}[\s)\-]?\s?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}$`

const maxEmailLength = 254

const minPwdLenght = 8
const maxPwdLenght = 64
const maxPhoneLenght = 20
const maxNameLenght = 40

const maxPhotoSize = 10 << 20

func ValidateEmail(email string) bool {
	if len(email) > maxEmailLength {
		return false
	}
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

func ValidatePhone(phone string) bool {
	if len(phone) > maxPhoneLenght {
		return false
	}
	re := regexp.MustCompile(phoneRegexp)
	return re.MatchString(phone)

}

func ValidateName(name string) bool {
	if len(name) > maxNameLenght {
		return false
	}
	return true
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

func ValidateProfile(phone string, firstname string, lastname string) error {
	if !ValidatePhone(phone) {
		return entity.NewValidationError("phone")
	}
	if !ValidateName(firstname) {
		return entity.NewValidationError("firstname")
	}
	if !ValidateName(lastname) {
		return entity.NewValidationError("lastname")
	}
	return nil
}

func ValidatePhoto(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}

	if fileHeader.Size <= 0 || fileHeader.Size > maxPhotoSize {
		return false
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".svg" {
		return false
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/svg+xml" {
		return false
	}

	return true
}
