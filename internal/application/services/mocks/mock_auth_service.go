// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/waliqueiroz/letmeask-api/internal/application/services (interfaces: AuthService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dtos "github.com/waliqueiroz/letmeask-api/internal/application/dtos"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthService) Login(arg0 dtos.CredentialsDTO) (dtos.AuthDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(dtos.AuthDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), arg0)
}
