package validator

import (
	"html"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func SanitizePosterInput(dto *dto.PosterInputFlatDTO) {
	if dto == nil {
		return
	}

	dto.Description = html.EscapeString(dto.Description)
	dto.CategoryAlias = html.EscapeString(dto.CategoryAlias)
	dto.Address = html.EscapeString(dto.Address)
	dto.City = html.EscapeString(dto.City)
	if dto.District != nil {
		escaped := html.EscapeString(*dto.District)
		dto.District = &escaped
	}

	if dto.Alias != nil {
		escaped := html.EscapeString(*dto.Alias)
		dto.Alias = &escaped
	}
	for i, feature := range dto.Features {
		dto.Features[i] = html.EscapeString(feature)
	}
}

func SanitizeUserProfile(dto *dto.CreateUserDTO) {
	if dto == nil {
		return
	}
	dto.Email = html.EscapeString(dto.Email)
	dto.FirstName = html.EscapeString(dto.FirstName)
	dto.LastName = html.EscapeString(dto.LastName)
	dto.Phone = html.EscapeString(dto.Phone)
}
