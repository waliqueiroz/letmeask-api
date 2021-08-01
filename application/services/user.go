package services

import (
	"github.com/waliqueiroz/letmeask-api/application/dtos"
	"github.com/waliqueiroz/letmeask-api/application/errors"
	"github.com/waliqueiroz/letmeask-api/application/providers"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type UserService interface {
	FindAll() ([]entities.User, error)
	Create(user entities.User) (entities.User, error)
	FindByID(userID string) (entities.User, error)
	Update(userID string, user entities.User) (entities.User, error)
	Delete(userID string) error
	UpdatePassword(userID string, password dtos.PasswordDTO) error
}

type userService struct {
	userRepository   repositories.UserRepository
	securityProvider providers.SecurityProvider
}

func NewUserService(userRepository repositories.UserRepository, securityProvider providers.SecurityProvider) *userService {
	return &userService{
		userRepository,
		securityProvider,
	}
}

func (service *userService) FindAll() ([]entities.User, error) {
	return service.userRepository.FindAll()
}

func (service *userService) Create(user entities.User) (entities.User, error) {
	hashedPassword, err := service.securityProvider.Hash(user.Password)
	if err != nil {
		return entities.User{}, err
	}

	user.Password = string(hashedPassword)

	return service.userRepository.Create(user)
}

func (service *userService) FindByID(userID string) (entities.User, error) {
	return service.userRepository.FindByID(userID)
}

func (service *userService) Update(userID string, user entities.User) (entities.User, error) {
	return service.userRepository.Update(userID, user)
}

func (service *userService) Delete(userID string) error {
	return service.userRepository.Delete(userID)
}

func (service *userService) UpdatePassword(userID string, password dtos.PasswordDTO) error {
	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if err := service.securityProvider.Verify(user.Password, password.Current); err != nil {
		return errors.NewUnauthorizedError("a operação falhou, revise os dados e tente novamente")
	}

	hashedPassword, err := service.securityProvider.Hash(password.New)
	if err != nil {
		return err
	}

	return service.userRepository.UpdatePassword(userID, string(hashedPassword))
}
