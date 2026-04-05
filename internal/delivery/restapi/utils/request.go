package utils

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

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
