package repositories

import "github.com/waliqueiroz/letmeask-api/domain/entities"

type UserRepository interface {
	FindAll() ([]entities.User, error)
	Create(user entities.User) (entities.User, error)
	FindById(userID string) (entities.User, error)
}
