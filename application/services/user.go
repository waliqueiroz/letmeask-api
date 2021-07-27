package services

import (
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
