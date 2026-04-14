package validator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
const pwdRegexPattern = `^[a-zA-Z\d!@#$%^&*\-]{8,}$`

const phoneRegexp = `^(\+7|8)\s?[\s(]?\d{3}[\s)\-]?\s?\d{3}[\s\-]?\d{2}[\s\-]?\d{2}$`

const maxEmailLength = 255

const minPwdLenght = 8
const maxPwdLenght = 64
const maxPhoneLenght = 20
const MaxNameLenght = 40

const maxPriceSize = 10000000

const maxPhotoSize = 10 << 20
const MaxPhotosLength = 12

const maxPosterDescriptionLength = 3000
const minAddressLength = 5
const maxAddressLength = 500
const maxDistrictLength = 100
const maxFloorCount = 100
const maxFacilityAliasLength = 50

func ValidateEmail(email string) bool {
	if len(email) > maxEmailLength {
		return false
	}
	re := regexp.MustCompile(emailRegexPattern)
	return re.MatchString(email)
}

func ValidatePhone(phone string) bool {
	re := regexp.MustCompile(phoneRegexp)
	return re.MatchString(phone)

}

func ValidateName(name string) bool {
	return len(name) <= MaxNameLenght
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

func ValidateAddress(address string) bool {
	address = strings.TrimSpace(address)

	if len(address) < minAddressLength || len(address) > maxAddressLength {
		return false
	}

	return true
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
	if poster.Price <= 0 || poster.Price > maxPriceSize {
		return entity.NewValidationError("price")
	}

	if len(poster.Description) > maxPosterDescriptionLength {
		return entity.NewValidationError("description")
	}

	if poster.Area <= 0 {
		return entity.NewValidationError("area")
	}

	if poster.GeoLat < -90 || poster.GeoLat > 90 {
		return entity.NewValidationError("geo_lat")
	}

	if poster.GeoLon < -180 || poster.GeoLon > 180 {
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
