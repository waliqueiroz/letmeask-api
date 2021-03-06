// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/waliqueiroz/letmeask-api/internal/application/providers (interfaces: Authenticator)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthenticator is a mock of Authenticator interface.
type MockAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticatorMockRecorder
}

// MockAuthenticatorMockRecorder is the mock recorder for MockAuthenticator.
type MockAuthenticatorMockRecorder struct {
	mock *MockAuthenticator
}

// NewMockAuthenticator creates a new mock instance.
func NewMockAuthenticator(ctrl *gomock.Controller) *MockAuthenticator {
	mock := &MockAuthenticator{ctrl: ctrl}
	mock.recorder = &MockAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticator) EXPECT() *MockAuthenticatorMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockAuthenticator) CreateToken(arg0 string, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockAuthenticatorMockRecorder) CreateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockAuthenticator)(nil).CreateToken), arg0, arg1)
}

// ExtractUserID mocks base method.
func (m *MockAuthenticator) ExtractUserID(arg0 interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExtractUserID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtractUserID indicates an expected call of ExtractUserID.
func (mr *MockAuthenticatorMockRecorder) ExtractUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractUserID", reflect.TypeOf((*MockAuthenticator)(nil).ExtractUserID), arg0)
}
