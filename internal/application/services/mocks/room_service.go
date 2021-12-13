package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
)

type RoomServiceMock struct {
	mock.Mock
}

func NewRoomServiceMock() *RoomServiceMock {
	return &RoomServiceMock{}
}

func (service *RoomServiceMock) Create(room entities.Room) (entities.Room, error) {
	args := service.Called(room)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) FindByID(roomID string) (entities.Room, error) {
	args := service.Called(roomID)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) EndRoom(userID string, roomID string) (entities.Room, error) {
	args := service.Called(userID, roomID)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) CreateQuestion(roomID string, question entities.Question) (entities.Room, error) {
	args := service.Called(roomID, question)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) UpdateQuestion(userID string, roomID string, questionID string, questionData dtos.UpdateQuestionDTO) (entities.Room, error) {
	args := service.Called(userID, roomID, questionID, questionData)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) LikeQuestion(roomID string, questionID string, like entities.Like) (entities.Room, error) {
	args := service.Called(roomID, questionID, like)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) DeslikeQuestion(roomID string, questionID string, likeID string) (entities.Room, error) {
	args := service.Called(roomID,
		questionID,
		likeID)

	return args.Get(0).(entities.Room), args.Error(1)
}

func (service *RoomServiceMock) DeleteQuestion(userID string, roomID string, questionID string) (entities.Room, error) {
	args := service.Called(userID, roomID, questionID)

	return args.Get(0).(entities.Room), args.Error(1)
}
