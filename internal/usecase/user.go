package usecase

import (
	"encoding/base64"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"mainService/internal/repository/redisTLC"

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

func (ucase *UserUsecase) UpdateUser(userID string, updInfo *domain.ApiUserUpdate) error {
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
	err := ucase.userRepo.UpdatePet(userID, petID, updInfo)
	if err != nil {
		return err
	}

	return nil
}
