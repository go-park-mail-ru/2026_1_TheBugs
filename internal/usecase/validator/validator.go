package validator

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const pwdRegexPattern = `^[a-zA-Z\d!@#$%^&*\-]{8,}$`

const phoneRegexp = `^(\+7|8)\s?[\s(]?\d{3}[\s)\-]?\s?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}$`
const addressRegexp = `^[а-яА-ЯёЁ\s\-\,\.\d\/]+$`

const maxEmailLength = 254

const minPwdLenght = 8
const maxPwdLenght = 64
const maxPhoneLenght = 20
const maxNameLenght = 40

const maxPhotoSize = 10 << 20
const MaxPhotosLength = 12

const maxPosterDescriptionLength = 499
const minAddressLength = 5
const maxAddressLength = 149
const maxDistrictLength = 29
const maxFloorCount = 99
const maxFacilityAliasLength = 49

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

func ValidatePhotos(photos []dto.PhotoInputDTO) error {
	if len(photos) == 0 {
		return entity.NewValidationError("min photos len")
	}
	if len(photos) > MaxPhotosLength {
		return entity.NewValidationError("max photos len")
	}
	for i, photo := range photos {
		if !ValidatePhoto(photo.FileHeader) {
			return entity.NewValidationError(fmt.Sprintf("photos[%d]", i))
		}
	}
	return nil
}

func ValidatePosterInputFlat(poster *dto.PosterInputFlatDTO) error {
	if poster.FlatNumber != nil && *poster.FlatNumber <= 0 {
		return entity.NewValidationError("flat_number")
	}

	if poster.FlatFloor <= 0 {
		return entity.NewValidationError("flat_floor")
	}

	if poster.FloorCount <= 0 {
		return entity.NewValidationError("floor_count")
	}

	if poster.FlatFloor > poster.FloorCount {
		return entity.NewValidationError("flat_floor")
	}

	return nil
}

func ValidateAddress(address string) bool {
	address = strings.TrimSpace(address)

	if len(address) < minAddressLength || len(address) > maxAddressLength {
		return false
	}

	re := regexp.MustCompile(addressRegexp)
	return re.MatchString(address)
}

func ValidateDistrict(district *string) bool {
	if district == nil {
		return true
	}

	return len(*district) <= maxDistrictLength
}

func ValidateFeatures(features []string) bool {
	for _, feature := range features {
		if feature == "" {
			return false
		}
		if len(feature) > maxFacilityAliasLength {
			return false
		}
	}

	return true
}

func ValidatePosterBase(poster *dto.PosterInputFlatDTO) error {
	if poster.Price <= 0 {
		return entity.NewValidationError("price")
	}

	if len(poster.Description) > maxPosterDescriptionLength {
		return entity.NewValidationError("description")
	}

	if poster.Area <= 0 {
		return entity.NewValidationError("area")
	}

	if poster.GeoLat < -90 && poster.GeoLat > 90 {
		return entity.NewValidationError("geo_lat")
	}

	if poster.GeoLon < -180 && poster.GeoLon > 180 {
		return entity.NewValidationError("geo_lon")
	}

	if !ValidateAddress(poster.Address) {
		return entity.NewValidationError("address")
	}

	if !ValidateDistrict(poster.District) {
		return entity.NewValidationError("district")
	}

	if poster.FloorCount <= 0 || poster.FloorCount > maxFloorCount {
		return entity.NewValidationError("floor_count")
	}

	if !ValidateFeatures(poster.Features) {
		return entity.NewValidationError("features")
	}

	return nil
}
