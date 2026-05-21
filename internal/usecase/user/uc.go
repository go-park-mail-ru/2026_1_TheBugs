package user

import (
	"context"
	"fmt"
	"html"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type UserUseCase struct {
	uow    usecase.UnitOfWork
	file   usecase.FileRepo
	sender usecase.MailSender
}

func NewUserUseCase(uow usecase.UnitOfWork, file usecase.FileRepo, sender usecase.MailSender) *UserUseCase {
	return &UserUseCase{
		uow:    uow,
		file:   file,
		sender: sender,
	}
}

func (uc *UserUseCase) GetByID(ctx context.Context, userID int) (*dto.UserDTO, error) {
	user, err := uc.uow.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}
	user.MakeAvatarPath()
	return user, nil
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, data dto.UpdateProfileRequest) (*dto.UserDTO, error) {
	if data.Phone != nil {
		if ok := validator.ValidatePhone(*data.Phone); !ok {
			return nil, entity.NewValidationError("phone")
		}
		*data.Phone = validator.NormolizePhoneNumber(*data.Phone)
	}
	if data.FirstName != nil {
		if ok := validator.ValidateName(*data.FirstName); !ok {
			return nil, entity.NewValidationError("first_name")
		}
		*data.FirstName = html.EscapeString(*data.FirstName)
	}
	if data.LastName != nil {
		if ok := validator.ValidateName(*data.LastName); !ok {
			return nil, entity.NewValidationError("last_name")
		}
		*data.LastName = html.EscapeString(*data.LastName)
	}
	updateDTO := dto.UpdateProfileDTO{ID: data.ID, FirstName: data.FirstName, LastName: data.LastName, Phone: data.Phone}
	if data.Avatar != nil {
		defer data.Avatar.File.Close()

		if !validator.ValidatePhoto(data.Avatar) {
			return nil, entity.NewValidationError("photo")
		}
		user, err := uc.uow.Users().GetByID(ctx, data.ID)
		if err != nil {
			return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
		}
		path := dto.GenerateAvatarPathForUser(data.ID)
		key := photo.GetKeyFromPath(path)
		size := data.Avatar.Size
		contentType := data.Avatar.ContentType

		if err := uc.file.Upload(ctx, key, data.Avatar.File, size, contentType); err != nil {
			return nil, fmt.Errorf(" uc.file.Upload: %w", err)
		}
		if user.AvatarURL != nil {
			if err := uc.file.Delete(ctx, photo.GetKeyFromPath(*user.AvatarURL)); err != nil {
				return nil, fmt.Errorf(" uc.file.Upload: %w", err)
			}
		}
		updateDTO.AvatarPath = &path
	}
	user, err := uc.uow.Users().UpdateProfile(ctx, updateDTO)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().UpdateProfile: %w", err)
	}
	user.MakeAvatarPath()
	return user, nil
}

func (uc *UserUseCase) GetRoommateUser(ctx context.Context, userID int) (*dto.RoommateUserProfileDTO, error) {
	user, err := uc.uow.Users().GetRoommateUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateUser: %w", err)
	}

	tags, err := uc.uow.Users().GetRoommateTags(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateTags: %w", err)
	}

	return dto.RoommateUserToDTO(user, tags), nil
}

func (uc *UserUseCase) AddRoommateMatch(ctx context.Context, fromUserID int, toUserID int) error {
	op := "UserUseCase.AddRoommateMatch"
	log := ctxLogger.GetLogger(ctx).WithField("op", op)

	if fromUserID == toUserID {
		return entity.InvalidInput
	}

	err := uc.uow.Users().AddRoommateMatch(ctx, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().AddRoommateMatch: %w", err)
	}

	toContacts, err := uc.uow.Users().GetRoommateContacts(ctx, toUserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().GetRoommateContacts: %w", err)
	}

	fromUser, err := uc.uow.Users().GetByID(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}

	poster, err := uc.uow.Posters().GetRoommatePoster(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().GetRoommatePoster: %w", err)
	}

	if uc.sender != nil {
		if err := uc.sender.SendRoommateMatch(
			ctx,
			toContacts.Email,
			fromUser.FirstName,
			fromUser.LastName,
			poster.Address,
		); err != nil {
			log.Errorf("uc.sender.SendRoommateMatch: %v", err)
		}
	}

	return nil
}

func (uc *UserUseCase) GetRoommateContacts(ctx context.Context, fromUserID int, toUserID int) (*dto.RoommateContactsDTO, error) {
	op := "UserUseCase.GetRoommateContacts"
	log := ctxLogger.GetLogger(ctx).WithField("op", op)

	if fromUserID == toUserID {
		return nil, entity.InvalidInput
	}

	isMatched, err := uc.uow.Users().IsRoommateMatch(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().IsRoommateMatch: %w", err)
	}

	if !isMatched {
		return nil, entity.NotFoundError
	}

	fromContacts, err := uc.uow.Users().GetRoommateContacts(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateContacts: %w", err)
	}

	toContacts, err := uc.uow.Users().GetRoommateContacts(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateContacts: %w", err)
	}

	fromUser, err := uc.uow.Users().GetByID(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}

	toUser, err := uc.uow.Users().GetByID(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}

	poster, err := uc.uow.Posters().GetRoommatePoster(ctx, fromUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetRoommatePoster: %w", err)
	}

	if uc.sender != nil {
		if err := uc.sender.SendRoommateContactsForRequester(
			ctx,
			fromContacts.Email,
			toUser.FirstName,
			toUser.LastName,
			toContacts.Email,
			toContacts.Phone,
			poster.Address,
		); err != nil {
			log.Errorf("uc.sender.SendRoommateContactsForRequester: %v", err)
		}

		if err := uc.sender.SendRoommateContactsForAccepted(
			ctx,
			toContacts.Email,
			fromUser.FirstName,
			fromUser.LastName,
			fromContacts.Email,
			fromContacts.Phone,
			poster.Address,
		); err != nil {
			log.Errorf("uc.sender.SendRoommateContactsForAccepted: %v", err)
		}
	}

	return toContacts, nil
}

func (uc *UserUseCase) CreateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error {
	if data.UserID <= 0 {
		return entity.InvalidInput
	}

	if data.Gender != "male" && data.Gender != "female" {
		return entity.InvalidInput
	}

	if data.Birthday == "" {
		return entity.InvalidInput
	}

	if data.Description == "" {
		return entity.InvalidInput
	}

	err := uc.uow.Users().CreateRoommateForm(ctx, data)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().CreateRoommateForm: %w", err)
	}

	return nil
}

func (uc *UserUseCase) GetRoommateForm(ctx context.Context, userID int) (*dto.RoommateFormDTO, error) {
	if userID <= 0 {
		return nil, entity.InvalidInput
	}

	form, err := uc.uow.Users().GetRoommateForm(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateForm: %w", err)
	}

	tags, err := uc.uow.Users().GetRoommateFormTags(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateFormTags: %w", err)
	}

	return dto.RoommateFormToDTO(form, tags), nil
}

func (uc *UserUseCase) UpdateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error {
	if data.UserID <= 0 {
		return entity.InvalidInput
	}

	if data.Gender != "male" && data.Gender != "female" {
		return entity.InvalidInput
	}

	if data.Birthday == "" {
		return entity.InvalidInput
	}

	if data.Description == "" {
		return entity.InvalidInput
	}

	err := uc.uow.Users().UpdateRoommateForm(ctx, data)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().UpdateRoommateForm: %w", err)
	}

	return nil
}

func (uc *UserUseCase) GetIncomingRoommateMatches(ctx context.Context, userID int) (*dto.RoommateMatchesResponse, error) {
	if userID <= 0 {
		return nil, entity.InvalidInput
	}

	users, err := uc.uow.Users().GetIncomingRoommateMatches(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetIncomingRoommateMatches: %w", err)
	}

	return &dto.RoommateMatchesResponse{
		Users: users,
		Len:   len(users),
	}, nil
}

func (uc *UserUseCase) GetMatchedRoommateMatches(ctx context.Context, userID int) (*dto.RoommateMatchesResponse, error) {
	if userID <= 0 {
		return nil, entity.InvalidInput
	}

	users, err := uc.uow.Users().GetMatchedRoommateMatches(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetMatchedRoommateMatches: %w", err)
	}

	return &dto.RoommateMatchesResponse{
		Users: users,
		Len:   len(users),
	}, nil
}
