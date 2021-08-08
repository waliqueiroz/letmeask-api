package services

import (
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type RoomService interface {
	Create(room entities.Room) (entities.Room, error)
	FindByID(roomID string) (entities.Room, error)
}

type roomService struct {
	roomRepository repositories.RoomRepository
}

func NewRoomService(roomRepository repositories.RoomRepository) *roomService {
	return &roomService{
		roomRepository,
	}
}

func (service *roomService) Create(room entities.Room) (entities.Room, error) {
	return service.roomRepository.Create(room)
}

func (service *roomService) FindByID(roomID string) (entities.Room, error) {
	return service.roomRepository.FindByID(roomID)
}
