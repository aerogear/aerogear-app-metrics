// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/web/interfaces.go

// Package mock_web is a generated GoMock package.
package mock_web

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPingable is a mock of Pingable interface
type MockPingable struct {
	ctrl     *gomock.Controller
	recorder *MockPingableMockRecorder
}

// MockPingableMockRecorder is the mock recorder for MockPingable
type MockPingableMockRecorder struct {
	mock *MockPingable
}

// NewMockPingable creates a new mock instance
func NewMockPingable(ctrl *gomock.Controller) *MockPingable {
	mock := &MockPingable{ctrl: ctrl}
	mock.recorder = &MockPingableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPingable) EXPECT() *MockPingableMockRecorder {
	return m.recorder
}

// Ping mocks base method
func (m *MockPingable) Ping() error {
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping
func (mr *MockPingableMockRecorder) Ping() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPingable)(nil).Ping))
}