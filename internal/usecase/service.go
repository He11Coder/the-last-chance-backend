package usecase

import (
	"mainService/internal/domain"
	"mainService/internal/repository/mongoTLC"
	"strings"
)

type IServiceUsecase interface {
	AddService(userID string, service *domain.ApiService) error
	GetServiceByID(serviceID string) (*domain.ApiService, error)
	GetUserServices(userID string) ([]*domain.ApiService, error)
	DeleteService(userID, serviceID string) error
	SearchServices(queryString string) ([]*domain.ApiService, error)
}

type ServiceUsecase struct {
	serviceRepo mongoTLC.IServiceRepository
	userRepo    mongoTLC.IUserRepository
}

func NewServiceUsecase(
	serviceRepository mongoTLC.IServiceRepository,
	userRepository mongoTLC.IUserRepository,
) IServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: serviceRepository,
		userRepo:    userRepository,
	}
}

func (ucase *ServiceUsecase) AddService(userID string, service *domain.ApiService) error {
	isRole := domain.IsRole(service.Type)
	if !isRole {
		return INVALID_ROLE
	}

	if service.Title == "" {
		return EMPTY_TITLE
	}

	err := ucase.serviceRepo.AddService(userID, service)
	if err != nil {
		return err
	}

	return nil
}

func (ucase *ServiceUsecase) GetServiceByID(serviceID string) (*domain.ApiService, error) {
	service, err := ucase.serviceRepo.GetServiceByID(serviceID)
	if err != nil {
		return nil, err
	}

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

	return services, nil
}

func (ucase *ServiceUsecase) DeleteService(userID, serviceID string) error {
	err := ucase.serviceRepo.DeleteService(userID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

func (ucase *ServiceUsecase) SearchServices(queryString string) ([]*domain.ApiService, error) {
	queryString = strings.TrimSpace(queryString)
	if queryString == "" {
		return nil, EMPTY_SEARCH_STRING
	}

	services, err := ucase.serviceRepo.SearchServices(queryString)
	if err != nil {
		return nil, err
	}

	return services, err
}
