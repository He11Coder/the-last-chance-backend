package usecase

import (
	"encoding/base64"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"mainService/internal/repository/redisTLC"

	"mainService/pkg/nsfwFilter"
	"mainService/pkg/serverErrors"
	"mainService/pkg/swearWordsDetector"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	Login(cred *domain.LoginCredentials) (*domain.LoginResponse, error)
	AddUser(newUser *domain.ApiUserInfo) (*domain.LoginResponse, error)
	UpdateUser(userID string, updInfo *domain.ApiUserUpdate) error
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetUserAvatar(userID string) (string, error)
	GetUserPets(userID string) (*domain.PetIDList, error)
	AddPet(userID string, petInfo *domain.ApiPetInfo) (*domain.ApiPetInfo, error)
	DeletePet(userID, petID string) error
	UpdatePet(userID, petID string, updInfo *domain.ApiPetUpdate) error
}

type UserUsecase struct {
	userRepo    mongoTLC.IUserRepository
	sessionRepo redisTLC.IAuthRepository
}

func NewUserUsecase(
	userRepository mongoTLC.IUserRepository,
	sessionRepository redisTLC.IAuthRepository,
) IUserUsecase {
	return &UserUsecase{
		userRepo:    userRepository,
		sessionRepo: sessionRepository,
	}
}

func (ucase *UserUsecase) ValidateImagesForNSFW(avatar, backImage string) error {
	imagesToValidate := []string{}

	if avatar != "" {
		imagesToValidate = append(imagesToValidate, avatar)
	}
	if backImage != "" {
		imagesToValidate = append(imagesToValidate, backImage)
	}

	if len(imagesToValidate) != 0 {
		results := nsfwFilter.RunInParallel(imagesToValidate...)

		if avatar != "" {
			userImageRes := results[0]
			if userImageRes.ProcessingErr != nil {
				return userImageRes.ProcessingErr
			}

			if !userImageRes.Inf.IsSafe {
				return serverErrors.NSFW_CONTENT_AVATAR_ERROR
			}

			if backImage != "" {
				backImageRes := results[1]
				if backImageRes.ProcessingErr != nil {
					return backImageRes.ProcessingErr
				}

				if !backImageRes.Inf.IsSafe {
					return serverErrors.NSFW_CONTENT_BACK_IMAGE_ERROR
				}
			}
		} else if backImage != "" {
			backImageRes := results[0]
			if backImageRes.ProcessingErr != nil {
				return backImageRes.ProcessingErr
			}

			if !backImageRes.Inf.IsSafe {
				return serverErrors.NSFW_CONTENT_BACK_IMAGE_ERROR
			}
		}
	}

	return nil
}

func (ucase *UserUsecase) Login(cred *domain.LoginCredentials) (*domain.LoginResponse, error) {
	userID, err := ucase.userRepo.CheckUser(cred)
	if err != nil {
		return nil, err
	}

	sessionID := uuid.NewString()

	err = ucase.sessionRepo.AddSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{UserID: userID, SessionID: sessionID}, nil
}

func (ucase *UserUsecase) AddUser(newUser *domain.ApiUserInfo) (*domain.LoginResponse, error) {
	validErr := ucase.ValidateImagesForNSFW(newUser.UserImage, newUser.UserBackImage)
	if validErr != nil {
		return nil, validErr
	}

	containsSwearWords := swearWordsDetector.DetectInMultipleInputs(newUser.Login, newUser.Username, newUser.Contacts)
	if containsSwearWords {
		return nil, serverErrors.SWEAR_WORDS_ERROR
	}

	verifStatus := ucase.userRepo.ValidateLogin(newUser.Login)
	if verifStatus != nil {
		return nil, verifStatus
	}

	if len(newUser.Password) == 0 {
		return nil, EMPTY_PASSWORD
	}

	userID, err := ucase.userRepo.AddUser(newUser)
	if err != nil {
		return nil, err
	}

	sessionID := uuid.NewString()

	err = ucase.sessionRepo.AddSession(sessionID, userID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{UserID: userID, SessionID: sessionID}, nil
}

func (ucase *UserUsecase) UpdateUser(userID string, updInfo *domain.ApiUserUpdate) error {
	validErr := ucase.ValidateImagesForNSFW(updInfo.UserImage, updInfo.UserBackImage)
	if validErr != nil {
		return validErr
	}

	containsSwearWords := swearWordsDetector.DetectInMultipleInputs(updInfo.Login, updInfo.Username, updInfo.Contacts)
	if containsSwearWords {
		return serverErrors.SWEAR_WORDS_ERROR
	}

	if updInfo.Login != "" {
		err := ucase.userRepo.ValidateLogin(updInfo.Login)
		if err != nil {
			return err
		}
	}

	if updInfo.NewPassword != "" {
		oldLogin, err := ucase.userRepo.GetUserLoginByID(userID)
		if err != nil {
			return err
		}

		oldCreds := &domain.LoginCredentials{
			Username: oldLogin,
			Password: updInfo.OldPassword,
		}

		_, err = ucase.userRepo.CheckUser(oldCreds)
		if err != nil {
			return err
		}
	}

	err := ucase.userRepo.UpdateUser(userID, updInfo)
	if err != nil {
		return err
	}

	return nil
}

func (ucase *UserUsecase) GetUserInfo(userID string) (*domain.ApiUserInfo, error) {
	uInfo, err := ucase.userRepo.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	uInfo.UserImage, err = ucase.GetUserAvatar(userID)
	if err != nil {
		return nil, err
	}

	return uInfo, nil
}

func (ucase *UserUsecase) GetUserAvatar(userID string) (string, error) {
	avaBytes, err := ucase.userRepo.GetAvatarBytes(userID)
	if err != nil {
		return "", nil
	}

	base64Image := ""
	if len(avaBytes) != 0 {
		base64Image = base64.StdEncoding.EncodeToString(avaBytes)
	}

	return base64Image, nil
}

func (ucase *UserUsecase) GetUserPets(userID string) (*domain.PetIDList, error) {
	petIDs, err := ucase.userRepo.GetUserPets(userID)
	if err != nil {
		return nil, err
	}

	apiPetIDs := &domain.PetIDList{
		PetIDs: petIDs,
	}

	return apiPetIDs, nil
}

func (ucase *UserUsecase) AddPet(userID string, petInfo *domain.ApiPetInfo) (*domain.ApiPetInfo, error) {
	validErr := ucase.ValidateImagesForNSFW(petInfo.PetAvatar, "")
	if validErr != nil {
		return nil, validErr
	}

	containsSwearWords := swearWordsDetector.DetectInMultipleInputs(petInfo.Info, petInfo.Name, petInfo.TypeOfAnimal)
	if containsSwearWords {
		return nil, serverErrors.SWEAR_WORDS_ERROR
	}

	petID, err := ucase.userRepo.AddPet(userID, petInfo)
	if err != nil {
		return nil, err
	}

	petIDStruct := &domain.ApiPetInfo{
		PetID: petID,
	}

	return petIDStruct, nil
}

func (ucase *UserUsecase) DeletePet(userID, petID string) error {
	err := ucase.userRepo.DeletePet(userID, petID)
	if err != nil {
		return err
	}

	return nil
}

func (ucase *UserUsecase) UpdatePet(userID, petID string, updInfo *domain.ApiPetUpdate) error {
	validErr := ucase.ValidateImagesForNSFW(updInfo.PetAvatar, "")
	if validErr != nil {
		return validErr
	}

	containsSwearWords := swearWordsDetector.DetectInMultipleInputs(updInfo.Info, updInfo.Name, updInfo.TypeOfAnimal)
	if containsSwearWords {
		return serverErrors.SWEAR_WORDS_ERROR
	}

	err := ucase.userRepo.UpdatePet(userID, petID, updInfo)
	if err != nil {
		return err
	}

	return nil
}
