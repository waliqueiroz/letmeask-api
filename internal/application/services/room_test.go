package services_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/waliqueiroz/letmeask-api/internal/application/dtos"
	application "github.com/waliqueiroz/letmeask-api/internal/application/errors"
	"github.com/waliqueiroz/letmeask-api/internal/application/services"
	"github.com/waliqueiroz/letmeask-api/internal/domain/entities"
	domain "github.com/waliqueiroz/letmeask-api/internal/domain/errors"
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

	Describe("Executing the FindByID function", func() {
		var roomID string
		var result entities.Room
		var findByIDError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, findByIDError = roomService.FindByID(roomID)
		})

		When("the FindByID function is executed with success", func() {
			var expectedFindByIDResult entities.Room

			BeforeEach(func() {
				roomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.FindByID result", func() {
				Expect(result).To(Equal(expectedFindByIDResult))
			})

			It("error should be nil", func() {
				Expect(findByIDError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(findByIDError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the EndRoom function", func() {
		var roomID string
		var userID string
		var result entities.Room
		var endRoomError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, endRoomError = roomService.EndRoom(userID, roomID)
		})

		When("the EndRoom function is executed with success", func() {
			var expectedEndRoomResult entities.Room

			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				endedRoomSerialized, err := ioutil.ReadFile("../../../test/resources/ended_room.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(endedRoomSerialized, &expectedEndRoomResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedEndRoomResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedEndRoomResult))
			})

			It("error should be nil", func() {
				Expect(endRoomError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(endRoomError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("the userID is not equal to the room author ID", func() {
			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				userID = "621f5e02e07fdbb81c8221l5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be a forbidden error", func() {
				Expect(endRoomError).To(Equal(application.NewForbiddenError("você não pode encerrar uma sala que não é sua.")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(endRoomError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the CreateQuestion function", func() {
		var roomID string
		var question entities.Question
		var result entities.Room
		var createQuestionError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, createQuestionError = roomService.CreateQuestion(roomID, question)
		})

		When("the CreateQuestion function is executed with success", func() {
			var expectedCreateQuestionResult entities.Room

			BeforeEach(func() {
				createQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedCreateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createQuestionRequestSerialized, &question)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedCreateQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedCreateQuestionResult))
			})

			It("error should be nil", func() {
				Expect(createQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				createQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createQuestionRequestSerialized, &question)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(createQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			var expectedUpdateResult entities.Room

			BeforeEach(func() {
				createQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/create_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedUpdateResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(createQuestionRequestSerialized, &question)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(createQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the UpdateQuestion function", func() {
		var userID string
		var roomID string
		var questionID string
		var questionData dtos.UpdateQuestionDTO
		var result entities.Room
		var updateQuestionError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, updateQuestionError = roomService.UpdateQuestion(userID, roomID, questionID, questionData)
		})

		When("the UpdateQuestion function is executed to highlight question with success", func() {
			var expectedUpdateQuestionResult entities.Room

			BeforeEach(func() {
				highlightQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/highlight_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionHighlightedSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_highlighted.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionHighlightedSerialized, &expectedUpdateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(highlightQuestionRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "621f5e02e07fdbb81c8221f5"
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedUpdateQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedUpdateQuestionResult))
			})

			It("error should be nil", func() {
				Expect(updateQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("the UpdateQuestion function is executed to mark question as answered with success", func() {
			var expectedUpdateQuestionResult entities.Room

			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionAnsweredSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_answered.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionAnsweredSerialized, &expectedUpdateQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(markQuestionAsAnsweredRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "621f5e02e07fdbb81c8221f5"
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedUpdateQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedUpdateQuestionResult))
			})

			It("error should be nil", func() {
				Expect(updateQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(markQuestionAsAnsweredRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				userID = "621f5e02e07fdbb81c8221f5"
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(updateQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("the userID is not equal to the room author ID", func() {
			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(markQuestionAsAnsweredRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "621f5e02e07fdbb81c8221e5"
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be a forbidden error", func() {
				Expect(updateQuestionError).To(Equal(application.NewForbiddenError("você não pode atualizar informações de uma sala que não é sua.")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			BeforeEach(func() {
				markQuestionAsAnsweredRequestSerialized, err := ioutil.ReadFile("../../../test/resources/mark_question_as_answered_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(markQuestionAsAnsweredRequestSerialized, &questionData)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				userID = "621f5e02e07fdbb81c8221f5"
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(updateQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the LikeQuestion function", func() {
		var roomID string
		var questionID string
		var like entities.Like
		var result entities.Room
		var likeQuestionError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, likeQuestionError = roomService.LikeQuestion(roomID, questionID, like)
		})

		When("the LikeQuestion function is executed with success", func() {
			var expectedLikeQuestionResult entities.Room

			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionLikedSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_liked.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionLikedSerialized, &expectedLikeQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedLikeQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedLikeQuestionResult))
			})

			It("error should be nil", func() {
				Expect(likeQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(likeQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("try to like question that does not exists", func() {
			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c82213f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should a ResourceNotFoundError", func() {
				Expect(likeQuestionError).To(Equal(domain.NewResourceNotFoundError("pergunta não encontrada.")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			BeforeEach(func() {
				likeQuestionRequestSerialized, err := ioutil.ReadFile("../../../test/resources/like_question_request.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(likeQuestionRequestSerialized, &like)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(likeQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the DeslikeQuestion function", func() {
		var roomID string
		var questionID string
		var likeID string
		var result entities.Room
		var deslikeQuestionError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, deslikeQuestionError = roomService.DeslikeQuestion(roomID, questionID, likeID)
		})

		When("the DeslikeQuestion function is executed with success", func() {
			var expectedDeslikeQuestionResult entities.Room

			BeforeEach(func() {
				roomWithQuestionLikedSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_liked.json")
				Expect(err).NotTo(HaveOccurred())

				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionLikedSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedDeslikeQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				likeID = "6176903081f4f3f262acd6b4"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedDeslikeQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedDeslikeQuestionResult))
			})

			It("error should be nil", func() {
				Expect(deslikeQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				likeID = "6176903081f4f3f262acd6b4"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(deslikeQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("try to deslike question that does not exists", func() {
			BeforeEach(func() {
				roomWithQuestionLikedSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_liked.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionLikedSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8421f9"
				likeID = "6176903081f4f3f262acd6b4"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should a ResourceNotFoundError", func() {
				Expect(deslikeQuestionError).To(Equal(domain.NewResourceNotFoundError("pergunta não encontrada.")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			BeforeEach(func() {
				roomWithQuestionLikedSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_question_liked.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionLikedSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				likeID = "6176903081f4f3f262acd6b4"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(deslikeQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})

	Describe("Executing the DeleteQuestion function", func() {
		var roomID string
		var userID string
		var questionID string
		var result entities.Room
		var deleteQuestionError error
		var roomService services.RoomService
		var mockCtrl *gomock.Controller

		JustBeforeEach(func() {
			result, deleteQuestionError = roomService.DeleteQuestion(userID, roomID, questionID)
		})

		When("the DeleteQuestion function is executed with success", func() {
			var expectedDeleteQuestionResult entities.Room

			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				roomSerialized, err := ioutil.ReadFile("../../../test/resources/room.json")
				Expect(err).NotTo(HaveOccurred())

				err = json.Unmarshal(roomSerialized, &expectedDeleteQuestionResult)
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(expectedDeleteQuestionResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be equal to expected roomRepository.Update result", func() {
				Expect(result).To(Equal(expectedDeleteQuestionResult))
			})

			It("error should be nil", func() {
				Expect(deleteQuestionError).Should(BeNil())
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while finding room by ID", func() {
			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.FindByID function", func() {
				Expect(deleteQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("the userID is not equal to the room author ID", func() {
			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID = "621f5e02e07fdbb81e8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be a forbidden error", func() {
				Expect(deleteQuestionError).To(Equal(application.NewForbiddenError("você não pode remover uma pergunta de uma sala que não é sua.")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})

		When("an error occurs while updating room in database", func() {
			BeforeEach(func() {
				roomWithQuestionsSerialized, err := ioutil.ReadFile("../../../test/resources/room_with_questions.json")
				Expect(err).NotTo(HaveOccurred())

				var expectedFindByIDResult entities.Room
				err = json.Unmarshal(roomWithQuestionsSerialized, &expectedFindByIDResult)
				Expect(err).NotTo(HaveOccurred())

				roomID = "621f5ec1e07fdbb81c8221f7"
				questionID = "621f5f94e07fdbb81c8221f9"
				userID = "621f5e02e07fdbb81c8221f5"

				mockCtrl = gomock.NewController(GinkgoT())

				mockRoomRepository := repositoriesMocks.NewMockRoomRepository(mockCtrl)
				mockRoomRepository.EXPECT().FindByID(roomID).Return(expectedFindByIDResult, nil).Times(1)
				mockRoomRepository.EXPECT().Update(roomID, gomock.AssignableToTypeOf(entities.Room{})).Return(entities.Room{}, errors.New("an error")).Times(1)

				roomService = services.NewRoomService(mockRoomRepository)
			})

			It("result should be an empty room struct", func() {
				Expect(result).To(Equal(entities.Room{}))
			})

			It("error should be the error returned by the roomRepository.Update function", func() {
				Expect(deleteQuestionError).To(Equal(errors.New("an error")))
			})

			AfterEach(func() {
				mockCtrl.Finish()
			})
		})
	})
})
