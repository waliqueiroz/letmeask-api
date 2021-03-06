package repositories

import "github.com/waliqueiroz/letmeask-api/internal/domain/entities"

type UserRepository interface {
	FindAll() ([]entities.User, error)
	Create(user entities.User) (entities.User, error)
	FindByID(userID string) (entities.User, error)
	Update(userID string, user entities.User) (entities.User, error)
	Delete(userID string) error
	UpdatePassword(userID string, password string) error
	FindByEmail(email string) (entities.User, error)
}
