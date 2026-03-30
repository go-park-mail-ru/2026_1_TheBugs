package dto

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type PosterCardDTO struct {
	ID           int     `json:"id"`
	Alias        string  `json:"alias"`
	Price        float64 `json:"price"`
	ImgURL       *string `json:"imageUrl"`
	Address      string  `json:"address"`
	Metro        *string `json:"metro"`
	Area         float64 `json:"area"`
	FlatCategory *string `json:"flat_category"`
	// Rating  *float64 `json:"rating"`
	// Beds    *int     `json:"beds"`
}

type PostersFiltersDTO struct {
	Limit          int
	Offset         int
	UtilityCompany *string
}

func PostersToPostersDTO(posters []entity.PosterFlat) []PosterCardDTO {
	listPosters := make([]PosterCardDTO, 0, len(posters))
	for _, poster := range posters {
		// rating := float64((i*3/2)%10 + 1)
		// count := (i%5 + 1)
		photoURL := poster.ImgURL
		if photoURL != nil {
			url := photo.MakeUrlFromPath(*photoURL, config.Config.Minio.PublicHost, config.Config.Minio.Bucket)
			photoURL = &url
		}
		posterDTO := PosterCardDTO{
			ID:      poster.ID,
			Price:   poster.Price,
			ImgURL:  photoURL,
			Address: poster.Address,
			Metro:   poster.Metro,
			Area:    poster.Area,
			// Rating:  &rating,
			// Beds:    &count,
			Alias:        poster.Alias,
			FlatCategory: poster.FlatCategory,
		}

		listPosters = append(listPosters, posterDTO)
	}
	return listPosters
}

type PosterDTO struct {
	ID          int             `json:"id"`
	Alias       string          `json:"alias"`
	Price       float64         `json:"price"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Area        float64         `json:"area"`
	Geo         GeographyDTO    `json:"building_geo"`
	Address     string          `json:"address"`
	District    *string         `json:"district"`
	Metro       *string         `json:"metro"`
	MetroGeo    *GeographyDTO   `json:"metro_geo"`
	City        string          `json:"city"`
	FloorCount  int             `json:"floor_count"`
	Images      []PhotoDTO      `json:"images"`
	Seller      PosterSellerDTO `json:"seller"`
	Flat        *FlatDTO        `json:"flat,omitempty"`
	House       *HouseDTO       `json:"house,omitempty"`
	Facilities  []FacilityDTO   `json:"facilities"`

	Company *UtilityCompanyCardDTO `json:"company,omitempty"`
}

func PosterToPosterDTO(poster *entity.PosterById) *PosterDTO {
	return &PosterDTO{
		ID:          poster.ID,
		Alias:       poster.Alias,
		Price:       poster.Price,
		Category:    poster.Category,
		Description: poster.Description,
		Area:        poster.Area,
		Geo:         GeographyPointToGeographyDTO(poster.Geo),
		Address:     poster.Address,
		District:    poster.District,
		Metro:       poster.Metro,
		MetroGeo:    GeographyPointPtrToGeographyDTO(poster.MetroGeo),
		City:        poster.City,
		FloorCount:  poster.FloorCount,
		Images:      posterImagesToPosterImagesDTO(poster.Images),
		Seller:      posterToPosterSellerDTO(poster),
		Company:     posterToUtilityCompanyCardDTO(poster),
		Facilities:  FacilitiesToFacilitiesDTO(poster.Facilities),
	}
}

// id: number;
//   alias: string;
//   address: string;
//   area: number;
//   price: number;
//   avatar_url: string;

type MyPosterDTO struct {
	ID        int     `json:"id"`
	Alias     string  `json:"alias"`
	Address   string  `json:"address"`
	Area      float64 `json:"area"`
	Price     float64 `json:"price"`
	AvatarURl *string `json:"avatar_url"`
}

func MyPosterToMyPosterDTO(posters []entity.Poster) []MyPosterDTO {
	listPosters := make([]MyPosterDTO, 0, len(posters))
	for _, poster := range posters {
		posterDTO := MyPosterDTO{
			ID:        poster.ID,
			Price:     poster.Price,
			AvatarURl: poster.AvatarURl,
			Address:   poster.Address,
			Area:      poster.Area,
			Alias:     poster.Alias,
		}

		listPosters = append(listPosters, posterDTO)
	}
	return listPosters
}

type PosterInputFlatDTO struct {
	UserID      int     `schema:"-"`
	Alias       *string `schema:"-"`
	Price       float64 `schema:"price"`
	Description string  `schema:"description"`
	CategoryID  int     `schema:"category_id"`
	Area        float64 `schema:"area"`

	GeoLat         float64 `schema:"geo_lat"`
	GeoLon         float64 `schema:"geo_lon"`
	FlatCategoryID int     `schema:"flat_category_id"`
	FlatNumber     *int    `schema:"flat_number"`
	FlatFloor      int     `schema:"flat_floor"`

	Address        string  `schema:"address"`
	CityID         int     `schema:"city_id"`
	MetroStationID *int    `schema:"metro_station_id"`
	District       *string `schema:"district"`
	FloorCount     int     `schema:"floor_count"`
	CompanyID      *int    `schema:"company_id"`

	Features []string        `schema:"features"`
	Images   []PhotoInputDTO `schema:"-"`
}

func PosterInputFlatDTOtoPosterInput(poster *PosterInputFlatDTO) *entity.PosterInput {
	return &entity.PosterInput{
		UserID:      poster.UserID,
		Price:       poster.Price,
		Description: poster.Description,

		CategoryID: poster.CategoryID,
		Area:       poster.Area,

		Address:        poster.Address,
		Geo:            GeographyInputDTOtoGeographyPoint(GeographyInputDTO{Lat: poster.GeoLat, Lon: poster.GeoLon}),
		CityID:         poster.CityID,
		MetroStationID: poster.MetroStationID,
		District:       poster.District,
		FloorCount:     poster.FloorCount,

		CompanyID: poster.CompanyID,

		Images: posterPhotosInputFlatDTOtoPhotosInput(poster),
	}
}

type CreatedPoster struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
}
