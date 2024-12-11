package usecase

import (
	"encoding/base64"
	"mainService/configs"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"mainService/internal/repository/redisTLC"
	"os"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	Login(cred *domain.LoginCredentials) (*domain.LoginResponse, error)
	AddUser(newUser *domain.ApiUserInfo) (*domain.LoginResponse, error)
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetUserAvatar(userID string) (string, error)
	GetUserPets(userID string) (*domain.PetIDList, error)
	AddPet(userID string, petInfo *domain.ApiPetInfo) error
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

	err := ucase.userRepo.AddPet(userID, petInfo)
	if err != nil {
		return err
	}

	return nil
}
