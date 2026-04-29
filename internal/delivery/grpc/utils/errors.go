package utils

import (
	"errors"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TranslateDomainsError(err error) error {
	switch {
	case errors.Is(err, entity.InvalidInput):

		var validateErr *entity.ValidationError
		if errors.As(err, &validateErr) {
			return status.Error(codes.InvalidArgument, validateErr.Error())
		}
		return status.Error(codes.InvalidArgument, "bad request")
	case errors.Is(err, entity.BadCredentials):
		return status.Error(codes.Unauthenticated, "unauthorized")
	case errors.Is(err, entity.NotFoundError):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, entity.AlredyExitError):
		return status.Error(codes.AlreadyExists, "record already existed")
	case errors.Is(err, entity.JWTError):
		return status.Error(codes.Unauthenticated, "bad jwt")
	case errors.Is(err, entity.OffsetOutOfRange):
		return status.Error(codes.NotFound, "not enough records")
	case errors.Is(err, entity.ToManyRequest):
		return status.Error(codes.ResourceExhausted, "too many requests")
	case errors.Is(err, entity.UnverifiedUser):
		return status.Error(codes.PermissionDenied, "email is unverified")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
