package validator

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const pwdRegexPattern = `^[a-zA-Z\d!@#$%^&*\-]{8,}$`

const phoneRegexp = `^(\+7|8)\s?[\s(]?\d{3}[\s)\-]?\s?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}$`

const maxEmailLength = 254

const minPwdLenght = 8
const maxPwdLenght = 64
const maxPhoneLenght = 20
const maxNameLenght = 40

const maxPriceSize = 10000000

const maxPhotoSize = 10 << 20
const MaxPhotosLength = 12

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
	return len(name) <= maxNameLenght
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

func ValidatePhoto(fileInput *dto.FileInput) bool {
	if fileInput == nil {
		return false
	}

	if fileInput.Size <= 0 || fileInput.Size > maxPhotoSize {
		return false
	}

	ext := filepath.Ext(fileInput.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".svg" {
		return false
	}

	contentType := fileInput.ContentType
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/svg+xml" {
		return false
	}

	return true
}

func ValidatePhotos(photos []dto.PhotoInputDTO) error {
	if len(photos) == 0 {
		return entity.NewValidationError("min photos len")
	}
	if len(photos) > MaxPhotosLength {
		return entity.NewValidationError("max photos len")
	}
	for i, photo := range photos {
		if photo.FileHeader == nil && photo.URL == nil {
			return entity.NewValidationError(fmt.Sprintf("photos[%d]", i))
		}
		if photo.FileHeader != nil && !ValidatePhoto(photo.FileHeader) {
			return entity.NewValidationError(fmt.Sprintf("photos[%d]", i))
		}
		if photo.Order < 0 {
			return entity.NewValidationError(fmt.Sprintf("photos[%d].order", i))
		}
		if photo.Order > MaxPhotosLength {
			return entity.NewValidationError(fmt.Sprintf("photos[%d].order", i))
		}
	}
	return nil
}

func ValidatePosterInputFlat(poster *dto.PosterInputFlatDTO) error {
	if poster.FlatFloor <= 0 || poster.FlatFloor > 500 {
		return entity.NewValidationError("flat_floor")
	}

	if poster.FloorCount <= 0 || poster.FloorCount > 500 {
		return entity.NewValidationError("floor_count")
	}

	if poster.FlatFloor > poster.FloorCount {
		return entity.NewValidationError("flat_floor")
	}
	if len(poster.Description) > 3000 {
		return entity.NewValidationError("descriprion")
	}
	if len(poster.Address) < 5 || len(poster.Address) > 500 {
		return entity.NewValidationError("address")
	}
	if poster.District != nil && len(*poster.District) > 100 {
		return entity.NewValidationError("district")
	}
	if poster.Price <= 0 || poster.Price > maxPriceSize {
		return entity.NewValidationError("price")
	}

	return nil
}
