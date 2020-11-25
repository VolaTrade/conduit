// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/tickers/internal/service (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/volatrade/tickers/internal/models"
	socket "github.com/volatrade/tickers/internal/socket"
	reflect "reflect"
	sync "sync"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// BuildOrderBookChannels mocks base method
func (m *MockService) BuildOrderBookChannels(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BuildOrderBookChannels", arg0)
}

// BuildOrderBookChannels indicates an expected call of BuildOrderBookChannels
func (mr *MockServiceMockRecorder) BuildOrderBookChannels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildOrderBookChannels", reflect.TypeOf((*MockService)(nil).BuildOrderBookChannels), arg0)
}

// BuildPairUrls mocks base method
func (m *MockService) BuildPairUrls() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildPairUrls")
	ret0, _ := ret[0].(error)
	return ret0
}

// BuildPairUrls indicates an expected call of BuildPairUrls
func (mr *MockServiceMockRecorder) BuildPairUrls() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildPairUrls", reflect.TypeOf((*MockService)(nil).BuildPairUrls))
}

// BuildTransactionChannels mocks base method
func (m *MockService) BuildTransactionChannels(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BuildTransactionChannels", arg0)
}

// BuildTransactionChannels indicates an expected call of BuildTransactionChannels
func (mr *MockServiceMockRecorder) BuildTransactionChannels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildTransactionChannels", reflect.TypeOf((*MockService)(nil).BuildTransactionChannels), arg0)
}

// CheckForDatabasePriveleges mocks base method
func (m *MockService) CheckForDatabasePriveleges(arg0 *sync.WaitGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CheckForDatabasePriveleges", arg0)
}

// CheckForDatabasePriveleges indicates an expected call of CheckForDatabasePriveleges
func (mr *MockServiceMockRecorder) CheckForDatabasePriveleges(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckForDatabasePriveleges", reflect.TypeOf((*MockService)(nil).CheckForDatabasePriveleges), arg0)
}

// GetOrderBookChannel mocks base method
func (m *MockService) GetOrderBookChannel(arg0 int) chan *models.OrderBookRow {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderBookChannel", arg0)
	ret0, _ := ret[0].(chan *models.OrderBookRow)
	return ret0
}

// GetOrderBookChannel indicates an expected call of GetOrderBookChannel
func (mr *MockServiceMockRecorder) GetOrderBookChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderBookChannel", reflect.TypeOf((*MockService)(nil).GetOrderBookChannel), arg0)
}

// GetSocketsArrayLength mocks base method
func (m *MockService) GetSocketsArrayLength() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSocketsArrayLength")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetSocketsArrayLength indicates an expected call of GetSocketsArrayLength
func (mr *MockServiceMockRecorder) GetSocketsArrayLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSocketsArrayLength", reflect.TypeOf((*MockService)(nil).GetSocketsArrayLength))
}

// GetTransactionChannel mocks base method
func (m *MockService) GetTransactionChannel(arg0 int) chan *models.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionChannel", arg0)
	ret0, _ := ret[0].(chan *models.Transaction)
	return ret0
}

// GetTransactionChannel indicates an expected call of GetTransactionChannel
func (mr *MockServiceMockRecorder) GetTransactionChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionChannel", reflect.TypeOf((*MockService)(nil).GetTransactionChannel), arg0)
}

// OrderBookChannelListenAndHandle mocks base method
func (m *MockService) OrderBookChannelListenAndHandle(arg0 chan *models.OrderBookRow, arg1 int, arg2 *sync.WaitGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OrderBookChannelListenAndHandle", arg0, arg1, arg2)
}

// OrderBookChannelListenAndHandle indicates an expected call of OrderBookChannelListenAndHandle
func (mr *MockServiceMockRecorder) OrderBookChannelListenAndHandle(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderBookChannelListenAndHandle", reflect.TypeOf((*MockService)(nil).OrderBookChannelListenAndHandle), arg0, arg1, arg2)
}

// ReportRunning mocks base method
func (m *MockService) ReportRunning(arg0 *sync.WaitGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportRunning", arg0)
}

// ReportRunning indicates an expected call of ReportRunning
func (mr *MockServiceMockRecorder) ReportRunning(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportRunning", reflect.TypeOf((*MockService)(nil).ReportRunning), arg0)
}

// SpawnSocketRoutines mocks base method
func (m *MockService) SpawnSocketRoutines(arg0 int) []*socket.BinanceSocket {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpawnSocketRoutines", arg0)
	ret0, _ := ret[0].([]*socket.BinanceSocket)
	return ret0
}

// SpawnSocketRoutines indicates an expected call of SpawnSocketRoutines
func (mr *MockServiceMockRecorder) SpawnSocketRoutines(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpawnSocketRoutines", reflect.TypeOf((*MockService)(nil).SpawnSocketRoutines), arg0)
}

// TransactionChannelListenAndHandle mocks base method
func (m *MockService) TransactionChannelListenAndHandle(arg0 chan *models.Transaction, arg1 int, arg2 *sync.WaitGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TransactionChannelListenAndHandle", arg0, arg1, arg2)
}

// TransactionChannelListenAndHandle indicates an expected call of TransactionChannelListenAndHandle
func (mr *MockServiceMockRecorder) TransactionChannelListenAndHandle(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionChannelListenAndHandle", reflect.TypeOf((*MockService)(nil).TransactionChannelListenAndHandle), arg0, arg1, arg2)
}
