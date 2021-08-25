package services

import (
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	"github.com/waliqueiroz/letmeask-api/internal/domain/repositories"
)

type RoomService interface {
	Create(room entities.Room) (entities.Room, error)
	FindByID(roomID string) (entities.Room, error)
	CreateQuestion(roomID string, question entities.Question) (entities.Room, error)
	LikeQuestion(roomID string, questionID string, like entities.Like) (entities.Room, error)
	DeslikeQuestion(roomID string, questionID string, likeID string) (entities.Room, error)
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

	room.AddQuestion(question)

	return service.roomRepository.Update(roomID, room)
}

func (service *roomService) LikeQuestion(roomID string, questionID string, like entities.Like) (entities.Room, error) {
	room, err := service.roomRepository.FindByID(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	err = room.LikeQuestion(questionID, like)
	if err != nil {
		return entities.Room{}, err
	}

	return service.roomRepository.Update(roomID, room)
}

func (service *roomService) DeslikeQuestion(roomID string, questionID string, likeID string) (entities.Room, error) {
	room, err := service.roomRepository.FindByID(roomID)
	if err != nil {
		return entities.Room{}, err
	}

	err = room.DeslikeQuestion(questionID, likeID)
	if err != nil {
		return entities.Room{}, err
	}

	return service.roomRepository.Update(roomID, room)
}
