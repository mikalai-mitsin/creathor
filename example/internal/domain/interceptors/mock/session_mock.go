// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/domain/interceptors (interfaces: SessionInterceptor)

// Package mock_interceptors is a generated GoMock package.
package mock_interceptors

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionInterceptor is a mock of SessionInterceptor interface.
type MockSessionInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockSessionInterceptorMockRecorder
}

// MockSessionInterceptorMockRecorder is the mock recorder for MockSessionInterceptor.
type MockSessionInterceptorMockRecorder struct {
	mock *MockSessionInterceptor
}

// NewMockSessionInterceptor creates a new mock instance.
func NewMockSessionInterceptor(ctrl *gomock.Controller) *MockSessionInterceptor {
	mock := &MockSessionInterceptor{ctrl: ctrl}
	mock.recorder = &MockSessionInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionInterceptor) EXPECT() *MockSessionInterceptorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionInterceptor) Create(arg0 context.Context, arg1 *models.SessionCreate, arg2 *models.User) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionInterceptorMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionInterceptor)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockSessionInterceptor) Delete(arg0 context.Context, arg1 string, arg2 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionInterceptorMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionInterceptor)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockSessionInterceptor) Get(arg0 context.Context, arg1 string, arg2 *models.User) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSessionInterceptorMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSessionInterceptor)(nil).Get), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockSessionInterceptor) List(arg0 context.Context, arg1 *models.SessionFilter, arg2 *models.User) ([]*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSessionInterceptorMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSessionInterceptor)(nil).List), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockSessionInterceptor) Update(arg0 context.Context, arg1 *models.SessionUpdate, arg2 *models.User) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSessionInterceptorMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessionInterceptor)(nil).Update), arg0, arg1, arg2)
}
