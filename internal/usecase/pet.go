package usecase

import (
	"encoding/base64"
	"mainService/configs"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"os"
)

type IPetUsecase interface {
	GetPetInfo(petID string) (*domain.ApiPetInfo, error)
	GetPetAvatar(petID string) (string, error)
}

type PetUsecase struct {
	petRepo mongoTLC.IPetRepository
}

func NewPetUsecase(
	petRepository mongoTLC.IPetRepository,
) IPetUsecase {
	return &PetUsecase{
		petRepo: petRepository,
	}
}

func (ucase *PetUsecase) GetPetInfo(petID string) (*domain.ApiPetInfo, error) {
	petInfo, err := ucase.GetPetInfo(petID)
	if err != nil {
		return nil, err
	}

	petInfo.PetAvatar, err = ucase.GetPetAvatar(petID)
	if err != nil {
		return nil, err
	}

	return petInfo, nil
}

func (ucase *PetUsecase) GetPetAvatar(petID string) (string, error) {
	avaPath, err := ucase.petRepo.GetAvatarPath(petID)
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
