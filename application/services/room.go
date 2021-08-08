package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/letmeask-api/domain/entities"
	"github.com/waliqueiroz/letmeask-api/domain/repositories"
)

type RoomService interface {
	Create(room entities.Room) (entities.Room, error)
	FindByID(roomID string) (entities.Room, error)
	CreateQuestion(roomID string, question entities.Question) (entities.Room, error)
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

func (service *roomService) CreateQuestion(roomID string, question entities.Question) (entities.Room, error) {
	room, err := service.roomRepository.FindByID(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	question.ID = uuid.New().String()
	question.CreatedAt = time.Now()
	room.Questions = append(room.Questions, question)

	return service.roomRepository.Update(roomID, room)
}
