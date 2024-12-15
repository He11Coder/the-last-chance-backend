package usecase

import (
	"encoding/base64"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
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
	petInfo, err := ucase.petRepo.GetPetInfo(petID)
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
	avaBytes, err := ucase.petRepo.GetAvatarBytes(petID)
	if err != nil {
		return "", nil
	}

	base64Image := ""
	if len(avaBytes) != 0 {
		base64Image = base64.StdEncoding.EncodeToString(avaBytes)
	}

	return base64Image, nil
}
