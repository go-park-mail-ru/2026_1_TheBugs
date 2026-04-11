package request

import "mime/multipart"

type UpdateProfileRequest struct {
	Phone     *string               `schema:"phone"`
	FirstName *string               `schema:"first_name"`
	LastName  *string               `schema:"last_name"`
	Avatar    *multipart.FileHeader `schema:"avatar"`
}
