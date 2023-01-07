// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/domain/usecases (interfaces: UserSessionUseCase)

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserSessionUseCase is a mock of UserSessionUseCase interface.
type MockUserSessionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserSessionUseCaseMockRecorder
}

// MockUserSessionUseCaseMockRecorder is the mock recorder for MockUserSessionUseCase.
type MockUserSessionUseCaseMockRecorder struct {
	mock *MockUserSessionUseCase
}

// NewMockUserSessionUseCase creates a new mock instance.
func NewMockUserSessionUseCase(ctrl *gomock.Controller) *MockUserSessionUseCase {
	mock := &MockUserSessionUseCase{ctrl: ctrl}
	mock.recorder = &MockUserSessionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserSessionUseCase) EXPECT() *MockUserSessionUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserSessionUseCase) Create(arg0 context.Context, arg1 *models.UserSessionCreate) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserSessionUseCaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserSessionUseCase)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockUserSessionUseCase) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserSessionUseCaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserSessionUseCase)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockUserSessionUseCase) Get(arg0 context.Context, arg1 string) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserSessionUseCaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserSessionUseCase)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockUserSessionUseCase) List(arg0 context.Context, arg1 *models.UserSessionFilter) ([]*models.UserSession, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.UserSession)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockUserSessionUseCaseMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUserSessionUseCase)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockUserSessionUseCase) Update(arg0 context.Context, arg1 *models.UserSessionUpdate) (*models.UserSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.UserSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserSessionUseCaseMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserSessionUseCase)(nil).Update), arg0, arg1)
}