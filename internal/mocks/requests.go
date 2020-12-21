// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/conduit/internal/requests (interfaces: Requests)

// Package requests is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRequests is a mock of Requests interface
type MockRequests struct {
	ctrl     *gomock.Controller
	recorder *MockRequestsMockRecorder
}

// MockRequestsMockRecorder is the mock recorder for MockRequests
type MockRequestsMockRecorder struct {
	mock *MockRequests
}

// NewMockRequests creates a new mock instance
func NewMockRequests(ctrl *gomock.Controller) *MockRequests {
	mock := &MockRequests{ctrl: ctrl}
	mock.recorder = &MockRequestsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequests) EXPECT() *MockRequestsMockRecorder {
	return m.recorder
}

// GetActiveBinanceExchangePairs mocks base method
func (m *MockRequests) GetActiveBinanceExchangePairs() ([]*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveBinanceExchangePairs")
	ret0, _ := ret[0].([]*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveBinanceExchangePairs indicates an expected call of GetActiveBinanceExchangePairs
func (mr *MockRequestsMockRecorder) GetActiveBinanceExchangePairs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveBinanceExchangePairs", reflect.TypeOf((*MockRequests)(nil).GetActiveBinanceExchangePairs))
}
