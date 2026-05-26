package user

import (
	"context"
	"fmt"
	"html"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
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

// FIXME: Переделать логику, добавить логику определения типа матча( входящий / исходящий), чтобы сообщения на почту отличались ( вас хотят добавить / вас приняли ...)
func (uc *UserUseCase) AddRoommateMatch(
	ctx context.Context,
	fromUserID int,
	toUserID int,
	posterAlias *string,
) error {
	op := "UserUseCase.AddRoommateMatch"
	log := ctxLogger.GetLogger(ctx).WithField("op", op)

	if fromUserID == toUserID {
		return entity.InvalidInput
	}

	var (
		isMatched bool

		fromContacts *dto.RoommateContactsDTO
		toContacts   *dto.RoommateContactsDTO

		fromUser *dto.UserDTO
		toUser   *dto.UserDTO

		poster *entity.PosterFlat
	)

	err := uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		_, err := r.Users().GetByID(ctx, toUserID)
		if err != nil {
			return fmt.Errorf("r.Users().GetByID: %w", err)
		}

		err = r.Users().AddRoommateMatch(ctx, fromUserID, toUserID, posterAlias)
		if err != nil {
			return fmt.Errorf("r.Users().AddRoommateMatch: %w", err)
		}

		isMatched, err = r.Users().IsRoommateMatch(ctx, fromUserID, toUserID)
		if err != nil {
			return fmt.Errorf("r.Users().IsRoommateMatch: %w", err)
		}

		toContacts, err = r.Users().GetRoommateContacts(ctx, toUserID)
		if err != nil {
			return fmt.Errorf("r.Users().GetRoommateContacts(to): %w", err)
		}

		fromUser, err = r.Users().GetByID(ctx, fromUserID)
		if err != nil {
			return fmt.Errorf("r.Users().GetByID(from): %w", err)
		}

		poster, err = r.Posters().GetRoommatePoster(ctx, fromUserID, toUserID)
		if err != nil {
			return fmt.Errorf("r.Posters().GetRoommatePoster: %w", err)
		}

		if !isMatched {
			return nil
		}

		fromContacts, err = r.Users().GetRoommateContacts(ctx, fromUserID)
		if err != nil {
			return fmt.Errorf("r.Users().GetRoommateContacts(from): %w", err)
		}

		toUser, err = r.Users().GetByID(ctx, toUserID)
		if err != nil {
			return fmt.Errorf("r.Users().GetByID(to): %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("uc.uow.Do: %w", err)
	}

	if uc.sender == nil {
		return nil
	}

	if !isMatched {
		if err := uc.sender.SendRoommateMatch(
			ctx,
			toContacts.Email,
			fromUser.FirstName,
			fromUser.LastName,
			poster.Alias,
		); err != nil {
			log.Errorf("uc.sender.SendRoommateMatch: %v", err)
		}

		return nil
	}

	if err := uc.sender.SendRoommateContactsForRequester(
		ctx,
		fromContacts.Email,
		toUser.FirstName,
		toUser.LastName,
		toContacts.Email,
		toContacts.Phone,
		poster.Alias,
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
		poster.Alias,
	); err != nil {
		log.Errorf("uc.sender.SendRoommateContactsForAccepted: %v", err)
	}

	return nil
}

func (uc *UserUseCase) DeleteRoommateMatch(ctx context.Context, fromUserID int, toUserID int) error {
	if fromUserID <= 0 || toUserID <= 0 {
		return entity.InvalidInput
	}

	if fromUserID == toUserID {
		return entity.InvalidInput
	}

	err := uc.uow.Users().DeleteRoommateMatch(ctx, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().DeleteRoommateMatch: %w", err)
	}

	return nil
}

func (uc *UserUseCase) GetRoommateContacts(ctx context.Context, fromUserID int, toUserID int) (*dto.RoommateContactsDTO, error) {
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

	toContacts, err := uc.uow.Users().GetRoommateContacts(ctx, toUserID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetRoommateContacts: %w", err)
	}

	return toContacts, nil
}

func (uc *UserUseCase) CreateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error {
	if data.UserID <= 0 {
		return entity.InvalidInput
	}

	if data.Gender != "male" && data.Gender != "female" {
		return &entity.ValidationError{Err: entity.InvalidInput, Field: "gender", Details: "only male or female"}
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

func (uc *UserUseCase) DeleteRoommateForm(ctx context.Context, userID int) error {
	if userID <= 0 {
		return entity.InvalidInput
	}

	err := uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		if err := r.Users().DeletePosterRoommatesByUserID(ctx, userID); err != nil {
			return fmt.Errorf("r.Users().DeletePosterRoommatesByUserID: %w", err)
		}

		if err := r.Users().DeleteRoommateForm(ctx, userID); err != nil {
			return fmt.Errorf("r.Users().DeleteRoommateForm: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("uc.uow.Do: %w", err)
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
	for i, user := range users {
		if user.AvatarURL != nil {
			avatar := photo.MakeUrlFromPath(*user.AvatarURL, config.Config.PublicHost, config.Config.Bucket)
			users[i].AvatarURL = &avatar
		}
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
	for i, user := range users {
		if user.AvatarURL != nil {
			avatar := photo.MakeUrlFromPath(*user.AvatarURL, config.Config.PublicHost, config.Config.Bucket)
			users[i].AvatarURL = &avatar
		}
	}

	return &dto.RoommateMatchesResponse{
		Users: users,
		Len:   len(users),
	}, nil
}
