package usecase

import (
	"encoding/base64"
	"mainService/configs"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"os"
)

type IUserUsecase interface {
	Login(cred *domain.LoginCredentials) (string, error)
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetUserAvatar(userID string) (string, error)
	GetUserPets(userID string) (*domain.PetIDList, error)
	AddPet(userID string, petInfo *domain.ApiPetInfo) error
}

type UserUsecase struct {
	userRepo mongoTLC.IUserRepository
}

func NewUserUsecase(
	userRepository mongoTLC.IUserRepository,
) IUserUsecase {
	return &UserUsecase{
		userRepo: userRepository,
	}
}

func (ucase *UserUsecase) Login(cred *domain.LoginCredentials) (string, error) {
	return "1", nil
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

	uInfo.UserImage = "1"

	return uInfo, nil
}

func (ucase *UserUsecase) GetUserAvatar(userID string) (string, error) {
	avaPath, err := ucase.userRepo.GetAvatarPath(userID)
	if err != nil {
		return "", nil
	}

	base64Image := ""
	if avaPath != "" {
		fileBytes, err := os.ReadFile(configs.CURR_DIR + avaPath)
		if err != nil {
			return "", err
		}

		base64Image = base64.StdEncoding.EncodeToString(fileBytes)
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

func (ucase *UserUsecase) AddPet(userID string, petInfo *domain.ApiPetInfo) error {
	petInfo.PetAvatar = ""

	err := ucase.AddPet(userID, petInfo)
	if err != nil {
		return err
	}

	return nil
}
