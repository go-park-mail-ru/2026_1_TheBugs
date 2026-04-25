package dto

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
	"github.com/google/uuid"
)

type OrderDTO struct {
	UserID      int             `schema:"-"`
	CategoryID  int             `schema:"category_id"`
	Description string          `schema:"description"`
	Images      []PhotoInputDTO `schema:"-"`
}

func OrderDTOtoOrder(order *OrderDTO) *Order {
	return &Order{
		UserID:      order.UserID,
		CategoryID:  order.CategoryID,
		Description: order.Description,
		Images:      orderPhotosDTOtoPhotos(order),
	}
}

type Order struct {
	UserID      int
	CategoryID  int
	Description string
	Images      []PhotoInput
}

func orderPhotosDTOtoPhotos(order *OrderDTO) []PhotoInput {
	photos := make([]PhotoInput, 0, len(order.Images))
	for _, photo := range order.Images {
		var photoInput PhotoInput
		photoInput.FileHeader = photo.FileHeader
		photoInput.Order = photo.Order
		photoInput.URL = photo.URL
		photos = append(photos, photoInput)
	}

	return photos
}

func MakePhotoPathsForOrder(order *Order) {
	for i, image := range order.Images {
		if image.FileHeader != nil {
			path := GeneratePhotoPathForOrder(image.Order)
			order.Images[i].Path = path
			continue
		}
		url := *image.URL
		order.Images[i].Path = strings.TrimPrefix(url, fmt.Sprintf("%s/%s", config.Config.PublicHost, config.Config.Bucket))
	}
}

func GeneratePhotoPathForOrder(order int) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	reqId := fmt.Sprintf("%016x", seed.Int())[:10]
	reqId += fmt.Sprintf("_ord_%d", order)
	return fmt.Sprintf("/poster/img/%s/%s.jpg", uuid.New().String(), reqId)
}

type OrderPreviewDTO struct {
	ID           int       `json:"id"`
	CategoryName string    `json:"category_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToOrderPreview(entities []entity.Order) []OrderPreviewDTO {
	result := make([]OrderPreviewDTO, 0, len(entities))

	for _, e := range entities {
		result = append(result, OrderPreviewDTO{
			ID:           e.ID,
			CategoryName: e.CategoryName,
			Status:       e.Status,
			CreatedAt:    e.CreatedAt,
		})
	}

	return result
}

type OrdersResponse struct {
	Len    int               `json:"len"`
	Orders []OrderPreviewDTO `json:"order"`
}

type OrderFullDTO struct {
	ID           int        `json:"id"`
	CategoryName string     `json:"category_name"`
	Status       string     `json:"status"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	Images       []PhotoDTO `json:"images"`
}

func OrderToOrderFullDTO(order *entity.OrderFull) *OrderFullDTO {
	images := make([]PhotoDTO, 0, len(order.Photos))

	for _, image := range order.Photos {
		images = append(images, PhotoDTO{
			ImgURL: image.ImgURL,
			Order:  image.Order,
		})
	}

	return &OrderFullDTO{
		ID:           order.ID,
		CategoryName: order.CategoryName,
		Status:       order.Status,
		Description:  order.Description,
		CreatedAt:    order.CreatedAt,
		Images:       images,
	}
}

func MakeOrderUrlsFromPaths(order *OrderFullDTO, publicHost string, bucket string) {
	for i, image := range order.Images {
		url := photo.MakeUrlFromPath(image.ImgURL, publicHost, bucket)
		order.Images[i].ImgURL = url
	}
}

type OrderAnswerDTO struct {
	Answer string `form:"answer"`
}
