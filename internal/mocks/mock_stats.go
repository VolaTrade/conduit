// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/conduit/internal/stats (interfaces: Stats)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStats is a mock of Stats interface
type MockStats struct {
	ctrl     *gomock.Controller
	recorder *MockStatsMockRecorder
}

// MockStatsMockRecorder is the mock recorder for MockStats
type MockStatsMockRecorder struct {
	mock *MockStats
}

// NewMockStats creates a new mock instance
func NewMockStats(ctrl *gomock.Controller) *MockStats {
	mock := &MockStats{ctrl: ctrl}
	mock.recorder = &MockStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStats) EXPECT() *MockStatsMockRecorder {
	return m.recorder
}

// Increment mocks base method
func (m *MockStats) Increment(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Increment", arg0)
}

// Increment indicates an expected call of Increment
func (mr *MockStatsMockRecorder) Increment(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockStats)(nil).Increment), arg0)
}

// ReportGoRoutines mocks base method
func (m *MockStats) ReportGoRoutines() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportGoRoutines")
}

// ReportGoRoutines indicates an expected call of ReportGoRoutines
func (mr *MockStatsMockRecorder) ReportGoRoutines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportGoRoutines", reflect.TypeOf((*MockStats)(nil).ReportGoRoutines))
}
