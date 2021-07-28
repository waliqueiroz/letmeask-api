package services

import (
	"errors"

	"github.com/waliqueiroz/letmeask-api/application/dtos"
	"github.com/waliqueiroz/letmeask-api/application/security"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository,
	}
}

func (service *UserService) FindAll() ([]entities.User, error) {
	return service.userRepository.FindAll()
}

func (service *UserService) Create(user entities.User) (entities.User, error) {
	hashedPassword, err := security.Hash(user.Password)
	if err != nil {
		return entities.User{}, err
	}

	user.Password = string(hashedPassword)

	return service.userRepository.Create(user)
}

func (service *UserService) FindByID(userID string) (entities.User, error) {
	return service.userRepository.FindByID(userID)
}

func (service *UserService) Update(userID string, user entities.User) (entities.User, error) {
	return service.userRepository.Update(userID, user)
}

func (service *UserService) Delete(userID string) error {
	return service.userRepository.Delete(userID)
}

func (service *UserService) UpdatePassword(userID string, password dtos.PasswordDTO) error {
	user, err := service.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if err := security.Verify(user.Password, password.Current); err != nil {
		return errors.New("a operação falhou. Revise os dados e tente novamente")
	}

	hashedPassword, err := security.Hash(password.New)
	if err != nil {
		return err
	}

	return service.userRepository.UpdatePassword(userID, string(hashedPassword))
}
