package dto

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type CategoryDTO struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

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

type PostersResponse struct {
	Len     int             `json:"len"`
	Posters []PosterCardDTO `json:"posters"`
}

type PostersFiltersDTO struct {
	Limit            int
	Offset           int
	SearchQuery      *string
	UtilityCompany   *string
	Category         *string
	MaxPrice         *int
	MinPrice         *int
	RoomCount        *int
	MaxSquare        *int
	MinSquare        *int
	Facilities       []string
	MaxFlatFloor     *int
	MinFlatFloor     *int
	IsNotFirstFloor  bool
	IsNotLastFloor   bool
	MaxBuildingFloor *int
	MinBuildingFloor *int
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
	Category    CategoryDTO     `json:"category"`
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
		Category:    CategoryDTO{Name: poster.Category, Alias: poster.CategoryAlias},
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

type MyPosterDTO struct {
	ID        int         `json:"id"`
	Alias     string      `json:"alias"`
	Address   string      `json:"address"`
	Area      float64     `json:"area"`
	Price     float64     `json:"price"`
	AvatarURl *string     `json:"avatar_url"`
	Category  CategoryDTO `json:"category"`
}

func MyPosterToMyPosterDTO(posters []entity.Poster) []MyPosterDTO {
	listPosters := make([]MyPosterDTO, 0, len(posters))
	for _, poster := range posters {
		photoURL := poster.AvatarURl
		if photoURL != nil {
			url := photo.MakeUrlFromPath(*photoURL, config.Config.Minio.PublicHost, config.Config.Minio.Bucket)
			photoURL = &url
		}
		posterDTO := MyPosterDTO{
			ID:        poster.ID,
			Price:     poster.Price,
			AvatarURl: photoURL,
			Address:   poster.Address,
			Area:      poster.Area,
			Alias:     poster.Alias,
			Category:  CategoryDTO{Name: poster.CategoryName, Alias: poster.CategoryAlias},
		}

		listPosters = append(listPosters, posterDTO)
	}
	return listPosters
}

type PosterInputFlatDTO struct {
	UserID        int     `schema:"-"`
	Alias         *string `schema:"-"`
	Price         float64 `schema:"price"`
	Description   string  `schema:"description"`
	CategoryAlias string  `schema:"category_alias"`
	Area          float64 `schema:"area"`

	GeoLat         float64 `schema:"geo_lat"`
	GeoLon         float64 `schema:"geo_lon"`
	FlatCategoryID int     `schema:"flat_category_id"`
	FlatNumber     *int    `schema:"flat_number"`
	FlatFloor      int     `schema:"flat_floor"`

	Address    string  `schema:"address"`
	City       string  `schema:"city"`
	District   *string `schema:"district"`
	FloorCount int     `schema:"floor_count"`
	CompanyID  *int    `schema:"company_id"`

	Features []string        `schema:"features"`
	Images   []PhotoInputDTO `schema:"-"`
}

type GenerateDescriptionDTO struct {
	Category     string   `json:"category"`
	Area         float64  `json:"area"`
	FlatCategory string   `json:"flat_category"`
	City         string   `json:"city"`
	Features     []string `json:"features"`
}

func PosterInputFlatDTOtoPosterInput(poster *PosterInputFlatDTO) *PosterInput {
	return &PosterInput{
		UserID:      poster.UserID,
		Price:       poster.Price,
		Description: poster.Description,

		CategoryAlias: poster.CategoryAlias,
		Area:          poster.Area,

		Address:    poster.Address,
		Geo:        GeographyInputDTOtoGeographyPoint(GeographyDTO{Lat: poster.GeoLat, Lon: poster.GeoLon}),
		District:   poster.District,
		FloorCount: poster.FloorCount,

		CompanyID: poster.CompanyID,

		Features: poster.Features,
		Images:   posterPhotosInputFlatDTOtoPhotosInput(poster),
	}
}

type CreatedPoster struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
}

type PosterInput struct {
	UserID int

	Alias       string
	Price       float64
	Description string

	CategoryAlias string
	Area          float64

	Address        string
	Geo            geo.GeographyPoint
	CityID         int
	MetroStationID *int
	District       *string
	FloorCount     int

	CompanyID *int

	Features []string
	Images   []PhotoInput
}

type PosterUpdateIDs struct {
	UserID     int
	PosterID   int
	PropertyID int
	BuildingID int
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type GeoJSONFeature struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"propertiese"`
	Geometry   Geometry       `json:"geometry"`
}

type GeoJSONFeatureResponse struct {
	Posters []GeoJSONFeature `json:"features"`
	Len     int              `json:"len"`
}

func PosterPointToGeoJSON(p entity.AnyPoint) GeoJSONFeature {
	properties := map[string]any{
		"id":    p.ID,
		"group": p.Group,
	}
	if p.Group {
		if p.Count != nil {
			properties["count"] = *p.Count
		}
		if p.PriceMin != nil {
			properties["priceMin"] = *p.PriceMin
		}
	} else {
		if p.Price != nil {
			properties["price"] = *p.Price
		}
		if p.Alias != nil {
			properties["alias"] = *p.Alias
		}
	}

	return GeoJSONFeature{
		Type:       "Feature",
		Properties: properties,
		Geometry: struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		}{
			Type:        "Point",
			Coordinates: []float64{p.Lon, p.Lat},
		},
	}
}

func PostersToGEOJsons(p []entity.AnyPoint) []GeoJSONFeature {
	res := make([]GeoJSONFeature, 0, len(p))
	for _, e := range p {
		res = append(res, PosterPointToGeoJSON(e))
	}
	return res
}

func ClustersToGEOJsons(clusters []entity.ClusterPoint) []GeoJSONFeature {
	out := make([]GeoJSONFeature, 0, len(clusters))

	for _, c := range clusters {
		out = append(out, GeoJSONFeature{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{c.Lon, c.Lat},
			},
			Properties: map[string]any{
				"id":      c.ID,
				"count":   c.Count,
				"cluster": true,
			},
		})
	}

	return out
}
