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
	GetUserAvatar(userID string)
	GetUserPets()
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

	if uInfo.UserImage != "" {
		fileBytes, err := os.ReadFile(configs.CURR_DIR + uInfo.Username)
		if err != nil {
			return nil, err
		}

		base64Image := base64.StdEncoding.EncodeToString(fileBytes)
		uInfo.UserImage = base64Image
	}

	return uInfo, nil
}

func (ucase *UserUsecase) GetUserAvatar(userID string) {
	_, err := ucase.userRepo.GetAvatarPath(userID)
	if err != nil {
		return
	}
}

func (ucase *UserUsecase) GetUserPets() {

}
