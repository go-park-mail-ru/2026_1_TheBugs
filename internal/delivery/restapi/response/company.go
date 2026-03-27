package response

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"

type DevelopersResponse struct {
	Len        int                `json:"len"`
	Developers []dto.DeveloperDTO `json:"developers"`
}

type CompaniesResponse struct {
	Len       int                         `json:"len"`
	Companies []dto.UtilityCompanyCardDTO `json:"utility_companies"`
}
