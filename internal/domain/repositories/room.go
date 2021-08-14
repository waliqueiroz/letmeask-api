package repositories

import "github.com/waliqueiroz/letmeask-api/internal/domain/entities"

type RoomRepository interface {
	Create(room entities.Room) (entities.Room, error)
	FindByID(roomID string) (entities.Room, error)
	Update(roomID string, room entities.Room) (entities.Room, error)
}
