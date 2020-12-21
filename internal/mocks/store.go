// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/conduit/internal/store (interfaces: StorageConnections)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/volatrade/conduit/internal/models"
	reflect "reflect"
)

// MockStorageConnections is a mock of StorageConnections interface
type MockStorageConnections struct {
	ctrl     *gomock.Controller
	recorder *MockStorageConnectionsMockRecorder
}

// MockStorageConnectionsMockRecorder is the mock recorder for MockStorageConnections
type MockStorageConnectionsMockRecorder struct {
	mock *MockStorageConnections
}

// NewMockStorageConnections creates a new mock instance
func NewMockStorageConnections(ctrl *gomock.Controller) *MockStorageConnections {
	mock := &MockStorageConnections{ctrl: ctrl}
	mock.recorder = &MockStorageConnectionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorageConnections) EXPECT() *MockStorageConnectionsMockRecorder {
	return m.recorder
}

// InsertOrderBookRowToDataBase mocks base method
func (m *MockStorageConnections) InsertOrderBookRowToDataBase(arg0 *models.OrderBookRow, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrderBookRowToDataBase", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrderBookRowToDataBase indicates an expected call of InsertOrderBookRowToDataBase
func (mr *MockStorageConnectionsMockRecorder) InsertOrderBookRowToDataBase(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrderBookRowToDataBase", reflect.TypeOf((*MockStorageConnections)(nil).InsertOrderBookRowToDataBase), arg0, arg1)
}

// InsertTransactionToDataBase mocks base method
func (m *MockStorageConnections) InsertTransactionToDataBase(arg0 *models.Transaction, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransactionToDataBase", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertTransactionToDataBase indicates an expected call of InsertTransactionToDataBase
func (mr *MockStorageConnectionsMockRecorder) InsertTransactionToDataBase(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransactionToDataBase", reflect.TypeOf((*MockStorageConnections)(nil).InsertTransactionToDataBase), arg0, arg1)
}

// MakeConnections mocks base method
func (m *MockStorageConnections) MakeConnections() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MakeConnections")
}

// MakeConnections indicates an expected call of MakeConnections
func (mr *MockStorageConnectionsMockRecorder) MakeConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeConnections", reflect.TypeOf((*MockStorageConnections)(nil).MakeConnections))
}

// TransferOrderBookCache mocks base method
func (m *MockStorageConnections) TransferOrderBookCache(arg0 []*models.OrderBookRow) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferOrderBookCache", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferOrderBookCache indicates an expected call of TransferOrderBookCache
func (mr *MockStorageConnectionsMockRecorder) TransferOrderBookCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferOrderBookCache", reflect.TypeOf((*MockStorageConnections)(nil).TransferOrderBookCache), arg0)
}

// TransferTransactionCache mocks base method
func (m *MockStorageConnections) TransferTransactionCache(arg0 []*models.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferTransactionCache", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferTransactionCache indicates an expected call of TransferTransactionCache
func (mr *MockStorageConnectionsMockRecorder) TransferTransactionCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferTransactionCache", reflect.TypeOf((*MockStorageConnections)(nil).TransferTransactionCache), arg0)
}