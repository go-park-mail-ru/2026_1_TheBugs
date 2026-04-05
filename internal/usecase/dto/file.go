package dto

import (
	"io"
)

type FileInput struct {
	Filename    string
	Size        int64
	ContentType string
	File        io.ReadCloser
}

/*
type PhotoInputtt struct {
	FileHeader *FileInput
	Path       string
	Order      int
}

type PhotoInputttDTO struct {
	FileHeader *FileInput
	Order      int
}

func (f *FileInput) Open() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(f.Data)), nil
}
*/
/* func ParsePhotos(r *http.Request) ([]PhotoInputttDTO, error) {
	var photos []PhotoInputttDTO
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			return nil, fmt.Errorf("r.ParseMultipartForm: %w", err)
		}
	}

	for i := 0; ; i++ {
		fileKey := fmt.Sprintf("photos.%d.file", i)
		orderKey := fmt.Sprintf("photos.%d.order", i)

		files := r.MultipartForm.File[fileKey]
		orders := r.MultipartForm.Value[orderKey]

		if len(files) == 0 && len(orders) == 0 {
			break
		}

		if len(files) == 0 {
			return nil, fmt.Errorf("photo %d: file is required", i)
		}
		if len(orders) == 0 {
			return nil, fmt.Errorf("photo %d: order is required", i)
		}

		order, err := strconv.Atoi(orders[0])
		if err != nil {
			return nil, fmt.Errorf("photo %d: invalid order: %w", i, err)
		}

		fileHeader := files[0]

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}

		photo := PhotoInputttDTO{
			FileHeader: &FileInput{
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

func ValidatePhoto(fileInput *FileInput) bool {
	if fileInput == nil {
		return false
	}

	maxPhotoSize := nil
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
*/
/* func ValidatePhotos(photos []PhotoInputttDTO) error {
	if len(photos) == 0 {
		return entity.NewValidationError("min photos len")
	}
	MaxPhotosLength := 0
	if len(photos) > MaxPhotosLength {
		return entity.NewValidationError("max photos len")
	}
	for i, photo := range photos {
		if !ValidatePhoto(photo.FileHeader) {
			return entity.NewValidationError(fmt.Sprintf("photos[%d]", i))
		}
		if photo.Order <= 0 {
			return entity.NewValidationError(fmt.Sprintf("photos[%d].order", i))
		}
		if photo.Order > MaxPhotosLength {
			return entity.NewValidationError(fmt.Sprintf("photos[%d].order", i))
		}
	}
	return nil
} */

/* func uploadPhoto(ctx context.Context, photoPoster PhotoInputtt) (string, error) {
	file, err := photoPoster.FileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("photoPoster.FileHeader.Open: %w", err)
	}
	defer file.Close()

	if ValidatePhoto(photoPoster.FileHeader) {
		return "", entity.NewValidationError("photo")
	}

	key := photo.GetKeyFromPath(photoPoster.Path)
	size := photoPoster.FileHeader.Size
	contentType := photoPoster.FileHeader.ContentType

	if err := uc.file.Upload(ctx, key, file, size, contentType); err != nil {
		return "", fmt.Errorf("uc.file.Upload: %w", err)
	}

	return key, nil
}
*/
