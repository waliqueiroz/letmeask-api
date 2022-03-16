package services_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	repositoriesMocks "github.com/waliqueiroz/letmeask-api/internal/infrastructure/database/mongodb/repositories/mocks"
)

var _ = Describe("Room", func() {

	Describe("Executing the Create function", func() {
		var room entities.Room
		var result entities.Room
		var createError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, createError = roomService.Create(room)
		})

		When("the Create function is executed with success", func() {
			var expectedCreateResult entities.Room

			BeforeEach(func() {
				createRoomRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_room_request.json")
				Expect(err).NotTo(HaveOccurred())

				createdRoomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createRoomRequestSerialized, &room)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createdRoomSerialized, &expectedCreateResult)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().Create(room).Return(expectedCreateResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Create result", func() {
				Expect(result).To(Equal(expectedCreateResult))
			})

			It("error should be nil", func() {
				Expect(createError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while saving room in database", func() {
			BeforeEach(func() {
				createRoomRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_room_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createRoomRequestSerialized, &room)
				Expect(err).NotTo(HaveOccurred())

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().Create(room).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Create function", func() {
				Expect(createError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

	})
})
