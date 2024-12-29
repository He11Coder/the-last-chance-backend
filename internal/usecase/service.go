package usecase

import (
	"encoding/base64"
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"mainService/pkg/serverErrors"
	"mainService/pkg/swearWordsDetector"
	"strings"
)

type IServiceUsecase interface {
	AddService(userID string, service *domain.ApiService) (*domain.ApiService, error)
	GetServiceByID(serviceID string) (*domain.ApiService, error)
	GetUserServices(userID string) ([]*domain.ApiService, error)
	GetAllServices() ([]*domain.ApiService, error)
	DeleteService(userID, serviceID string) error
	SearchServices(queryString string, filters *domain.ServiceFilter) ([]*domain.ApiService, error)
}

type ServiceUsecase struct {
	serviceRepo mongoTLC.IServiceRepository
	userRepo    mongoTLC.IUserRepository
	petRepo     mongoTLC.IPetRepository
}

func NewServiceUsecase(
	serviceRepository mongoTLC.IServiceRepository,
	userRepository mongoTLC.IUserRepository,
	petRepository mongoTLC.IPetRepository,
) IServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepository,
		userRepo:    userRepository,
		petRepo:     petRepository,
	}
}

func (ucase *ServiceUsecase) AddService(userID string, service *domain.ApiService) (*domain.ApiService, error) {
	isRole := domain.IsRole(service.Type)
	if !isRole {
		return nil, INVALID_ROLE
	}

	containsSwearWords := swearWordsDetector.DetectInMultipleInputs(service.Description, service.Title)
	if containsSwearWords {
		return nil, serverErrors.SWEAR_WORDS_ERROR
	}

	if service.Title == "" {
		return nil, EMPTY_TITLE
	}

	serviceID, err := ucase.serviceRepo.AddService(userID, service)
	if err != nil {
		return nil, err
	}

	serviceIDStruct := &domain.ApiService{
		ServiceID: serviceID,
	}

	if len(service.PetIDs) != 0 {
		for _, petID := range service.PetIDs {
			info, err := ucase.petRepo.GetPetInfo(petID)
			if err != nil {
				return nil, err
			}

			if info.TypeOfAnimal != "" {
				err = ucase.petRepo.IncrementAnimal(info.TypeOfAnimal, serviceID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return serviceIDStruct, nil
}

func (ucase *ServiceUsecase) GetServiceByID(serviceID string) (*domain.ApiService, error) {
	service, err := ucase.serviceRepo.GetServiceByID(serviceID)
	if err != nil {
		return nil, err
	}

	if len(service.PetIDs) == 0 {
		service.PetIDs = []string{}
	}

	avatar, err := ucase.userRepo.GetAvatarBytes(service.UserID)
	if err != nil {
		return nil, err
	}

	service.UserImage = base64.StdEncoding.EncodeToString(avatar)

	return service, nil
}

func (ucase *ServiceUsecase) GetUserServices(userID string) ([]*domain.ApiService, error) {
	serviceIDs, err := ucase.userRepo.GetUserServices(userID)
	if err != nil {
		return nil, err
	}

	services, err := ucase.serviceRepo.GetServicesByIDs(serviceIDs...)
	if err != nil {
		return nil, err
	}

	for _, serv := range services {
		if len(serv.PetIDs) == 0 {
			serv.PetIDs = []string{}
		}

		avatar, err := ucase.userRepo.GetAvatarBytes(serv.UserID)
		if err != nil {
			return nil, err
		}

		serv.UserImage = base64.StdEncoding.EncodeToString(avatar)
	}

	return services, nil
}

func (ucase *ServiceUsecase) GetAllServices() ([]*domain.ApiService, error) {
	services, err := ucase.serviceRepo.GetAllServices()
	if err != nil {
		return nil, err
	}

	for _, serv := range services {
		if len(serv.PetIDs) == 0 {
			serv.PetIDs = []string{}
		}

		avatar, err := ucase.userRepo.GetAvatarBytes(serv.UserID)
		if err != nil {
			return nil, err
		}

		serv.UserImage = base64.StdEncoding.EncodeToString(avatar)
	}

	return services, nil
}

func (ucase *ServiceUsecase) DeleteService(userID, serviceID string) error {
	servInfo, err := ucase.serviceRepo.GetServiceByID(serviceID)
	if err != nil {
		return err
	}

	if len(servInfo.PetIDs) != 0 {
		for _, petID := range servInfo.PetIDs {
			petInfo, err := ucase.petRepo.GetPetInfo(petID)
			if err != nil {
				return err
			}

			if petInfo.TypeOfAnimal != "" {
				err = ucase.petRepo.DecrementAnimal(petInfo.TypeOfAnimal, serviceID)
				if err != nil {
					return err
				}
			}
		}
	}

	err = ucase.serviceRepo.DeleteService(userID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

func (ucase *ServiceUsecase) SearchServices(queryString string, filters *domain.ServiceFilter) ([]*domain.ApiService, error) {
	if (filters.MinPrice > filters.MaxPrice && (filters.MaxPrice != 0)) || (filters.MinPrice < 0) || (filters.MaxPrice < 0) {
		return nil, INVALID_PRICE_RANGE
	}

	services, err := ucase.serviceRepo.SearchServices(strings.TrimSpace(queryString), filters)
	if err != nil {
		return nil, err
	}

	for _, serv := range services {
		if len(serv.PetIDs) == 0 {
			serv.PetIDs = []string{}
		}

		avatar, err := ucase.userRepo.GetAvatarBytes(serv.UserID)
		if err != nil {
			return nil, err
		}

		serv.UserImage = base64.StdEncoding.EncodeToString(avatar)
	}

	return services, err
}
