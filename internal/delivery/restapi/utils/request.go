package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func ParseIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	if idStr == "" {
		return 0, fmt.Errorf("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return id, nil
}

func ParseAliasFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	alias := vars["alias"]
	if alias == "" {
		return "", fmt.Errorf("alias is required")
	}

	return alias, nil
}

func ParseFormData(r *http.Request, form interface{}) error {
	var decoder = schema.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		log.Printf("r.ParseForm: %s", err)
		return fmt.Errorf("r.ParseForm, %w", err)
	}
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(form, r.PostForm)
	return err
}

func ParseMultipartFormData(r *http.Request, form interface{}) error {
	var decoder = schema.NewDecoder()
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		return fmt.Errorf("r.ParseMultipartForm, %w", err)
	}
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(form, r.PostForm)
	return err
}

func ParsePhotos(r *http.Request) ([]dto.PhotoInputDTO, error) {
	var photos []dto.PhotoInputDTO
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			return nil, fmt.Errorf("r.ParseMultipartForm: %w", err)
		}
	}

	for i := 0; ; i++ {
		fileKey := fmt.Sprintf("photos.%d.file", i)
		orderKey := fmt.Sprintf("photos.%d.order", i)
		urlKey := fmt.Sprintf("photos.%d.url", i)

		files := r.MultipartForm.File[fileKey]
		orders := r.MultipartForm.Value[orderKey]
		urls := r.MultipartForm.Value[urlKey]

		log.Printf("len(fiels): %d, order: %v, urls: %v", len(files), orders, urls)

		if len(files) == 0 && len(orders) == 0 {
			break
		}

		if len(files) == 0 && len(urls) == 0 {
			return nil, fmt.Errorf("photo %d: file or url is required", i)
		}
		if len(orders) == 0 {
			return nil, fmt.Errorf("photo %d: order is required", i)
		}

		order, err := strconv.Atoi(orders[0])
		if err != nil {
			return nil, fmt.Errorf("photo %d: invalid order: %w", i, err)
		}
		if len(files) == 0 {
			photo := dto.PhotoInputDTO{
				Order: order,
				URL:   &urls[0],
			}

			photos = append(photos, photo)
			continue
		}

		fileHeader := files[0]

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		photo := dto.PhotoInputDTO{
			FileHeader: &dto.FileInput{
				Filename:    fileHeader.Filename,
				Size:        fileHeader.Size,
				ContentType: fileHeader.Header.Get("Content-Type"),
				File:        file,
			},
			Order: order,
		}

		photos = append(photos, photo)
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].Order < photos[j].Order
	})

	return photos, nil
}

const (
	defaultLimit  = "12"
	defaultOffset = "0"
)

func ParsePostersFilters(r *http.Request) (dto.PostersFiltersDTO, error) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit == 0 {
		limit, _ = strconv.Atoi(defaultLimit)
	}

	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset == 0 {
		offset, _ = strconv.Atoi(defaultOffset)
	}

	params := dto.PostersFiltersDTO{
		Offset: offset,
		Limit:  limit,
	}
	if v := q.Get("facilities"); v != "" {
		params.Facilities = strings.Split(v, ",")
	}

	if v := q.Get("search_query"); v != "" {
		params.SearchQuery = &v
	}
	if v := q.Get("utility_company"); v != "" {
		params.UtilityCompany = &v
	}
	if v := q.Get("category"); v != "" {
		params.Category = &v
	}

	if v := q.Get("room_count"); v != "" {
		i, _ := strconv.Atoi(v)
		params.RoomCount = &i
	}
	setIntIfNotEmpty(q, "min_price", &params.MinPrice)
	setIntIfNotEmpty(q, "max_price", &params.MaxPrice)
	setIntIfNotEmpty(q, "min_square", &params.MinSquare)
	setIntIfNotEmpty(q, "max_square", &params.MaxSquare)
	setIntIfNotEmpty(q, "min_flat_floor", &params.MinFlatFloor)
	setIntIfNotEmpty(q, "max_flat_floor", &params.MaxFlatFloor)
	setIntIfNotEmpty(q, "min_building_floor", &params.MinBuildingFloor)
	setIntIfNotEmpty(q, "max_building_floor", &params.MaxBuildingFloor)

	params.IsNotFirstFloor = q.Get("not_first_floor") == "true"
	params.IsNotLastFloor = q.Get("not_last_floor") == "true"

	return params, nil
}

func setIntIfNotEmpty(q url.Values, key string, ptr **int) {
	if v := q.Get(key); v != "" {
		i, _ := strconv.Atoi(v)
		*ptr = &i
	}
}

func ParseFileInput(fileHeader *multipart.FileHeader) (*dto.FileInput, error) {
	if fileHeader == nil {
		return nil, nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("fileHeader.Open: %w", err)
	}

	return &dto.FileInput{
		Filename:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
		File:        file,
	}, nil
}

func ParseMapFilters(r *http.Request) (dto.MapBounds, error) {
	q := r.URL.Query()

	swLat, err := strconv.ParseFloat(q.Get("sw_lat"), 32)
	if err != nil {
		return dto.MapBounds{}, err
	}
	swLng, err := strconv.ParseFloat(q.Get("sw_lon"), 32)
	if err != nil {
		return dto.MapBounds{}, err
	}
	neLat, err := strconv.ParseFloat(q.Get("ne_lat"), 32)
	if err != nil {
		return dto.MapBounds{}, err
	}
	neLng, err := strconv.ParseFloat(q.Get("ne_lon"), 32)
	if err != nil {
		return dto.MapBounds{}, err
	}

	zoom, err := strconv.Atoi(q.Get("zoom"))
	if err != nil {
		return dto.MapBounds{}, err
	}
	return dto.MapBounds{
		BBox: dto.BBox{
			SouthWest: dto.SouthWest{
				Lat: swLat,
				Lon: swLng,
			},
			NorthEast: dto.NorthEast{
				Lat: neLat,
				Lon: neLng,
			},
		},
		Zoom: zoom,
	}, nil
}
