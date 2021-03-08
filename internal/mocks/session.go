// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/conduit/internal/session (interfaces: Session)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSession is a mock of Session interface.
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

// MockSessionMockRecorder is the mock recorder for MockSession.
type MockSessionMockRecorder struct {
	mock *MockSession
}

// NewMockSession creates a new mock instance.
func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

// GetConnectionCount mocks base method.
func (m *MockSession) GetConnectionCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetConnectionCount indicates an expected call of GetConnectionCount.
func (mr *MockSessionMockRecorder) GetConnectionCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionCount", reflect.TypeOf((*MockSession)(nil).GetConnectionCount))
}

// ReportRunning mocks base method.
func (m *MockSession) ReportRunning(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportRunning", arg0)
}

// ReportRunning indicates an expected call of ReportRunning.
func (mr *MockSessionMockRecorder) ReportRunning(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportRunning", reflect.TypeOf((*MockSession)(nil).ReportRunning), arg0)
}
