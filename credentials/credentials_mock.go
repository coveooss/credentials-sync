// Code generated by MockGen. DO NOT EDIT.
// Source: credentials/credentials.go

// Package credentials is a generated GoMock package.
package credentials

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCredentials is a mock of Credentials interface
type MockCredentials struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialsMockRecorder
}

// MockCredentialsMockRecorder is the mock recorder for MockCredentials
type MockCredentialsMockRecorder struct {
	mock *MockCredentials
}

// NewMockCredentials creates a new mock instance
func NewMockCredentials(ctrl *gomock.Controller) *MockCredentials {
	mock := &MockCredentials{ctrl: ctrl}
	mock.recorder = &MockCredentialsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCredentials) EXPECT() *MockCredentialsMockRecorder {
	return m.recorder
}

// BaseValidate mocks base method
func (m *MockCredentials) BaseValidate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseValidate")
	ret0, _ := ret[0].(error)
	return ret0
}

// BaseValidate indicates an expected call of BaseValidate
func (mr *MockCredentialsMockRecorder) BaseValidate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseValidate", reflect.TypeOf((*MockCredentials)(nil).BaseValidate))
}

// GetID mocks base method
func (m *MockCredentials) GetID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetID indicates an expected call of GetID
func (mr *MockCredentialsMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockCredentials)(nil).GetID))
}

// GetTargetID mocks base method
func (m *MockCredentials) GetTargetID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTargetID indicates an expected call of GetTargetID
func (mr *MockCredentialsMockRecorder) GetTargetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetID", reflect.TypeOf((*MockCredentials)(nil).GetTargetID))
}

// ShouldSync mocks base method
func (m *MockCredentials) ShouldSync(targetName string, targetTags map[string]string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldSync", targetName, targetTags)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ShouldSync indicates an expected call of ShouldSync
func (mr *MockCredentialsMockRecorder) ShouldSync(targetName, targetTags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldSync", reflect.TypeOf((*MockCredentials)(nil).ShouldSync), targetName, targetTags)
}

// ToString mocks base method
func (m *MockCredentials) ToString(arg0 bool) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToString", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// ToString indicates an expected call of ToString
func (mr *MockCredentialsMockRecorder) ToString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToString", reflect.TypeOf((*MockCredentials)(nil).ToString), arg0)
}

// Validate mocks base method
func (m *MockCredentials) Validate() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockCredentialsMockRecorder) Validate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockCredentials)(nil).Validate))
}
