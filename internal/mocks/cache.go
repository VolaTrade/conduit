// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/volatrade/conduit/internal/cache (interfaces: Cache)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/volatrade/conduit/internal/models"
	reflect "reflect"
)

// MockCache is a mock of Cache interface
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// GetAllOrderBookRows mocks base method
func (m *MockCache) GetAllOrderBookRows() []*models.OrderBookRow {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOrderBookRows")
	ret0, _ := ret[0].([]*models.OrderBookRow)
	return ret0
}

// GetAllOrderBookRows indicates an expected call of GetAllOrderBookRows
func (mr *MockCacheMockRecorder) GetAllOrderBookRows() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOrderBookRows", reflect.TypeOf((*MockCache)(nil).GetAllOrderBookRows))
}

// GetEntries mocks base method
func (m *MockCache) GetEntries() []*models.CacheEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntries")
	ret0, _ := ret[0].([]*models.CacheEntry)
	return ret0
}

// GetEntries indicates an expected call of GetEntries
func (mr *MockCacheMockRecorder) GetEntries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntries", reflect.TypeOf((*MockCache)(nil).GetEntries))
}

// InsertEntry mocks base method
func (m *MockCache) InsertEntry(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InsertEntry", arg0)
}

// InsertEntry indicates an expected call of InsertEntry
func (mr *MockCacheMockRecorder) InsertEntry(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertEntry", reflect.TypeOf((*MockCache)(nil).InsertEntry), arg0)
}

// InsertOrderBookRow mocks base method
func (m *MockCache) InsertOrderBookRow(arg0 *models.OrderBookRow) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InsertOrderBookRow", arg0)
}

// InsertOrderBookRow indicates an expected call of InsertOrderBookRow
func (mr *MockCacheMockRecorder) InsertOrderBookRow(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrderBookRow", reflect.TypeOf((*MockCache)(nil).InsertOrderBookRow), arg0)
}

// OrderBookRowsLength mocks base method
func (m *MockCache) OrderBookRowsLength() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderBookRowsLength")
	ret0, _ := ret[0].(int)
	return ret0
}

// OrderBookRowsLength indicates an expected call of OrderBookRowsLength
func (mr *MockCacheMockRecorder) OrderBookRowsLength() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderBookRowsLength", reflect.TypeOf((*MockCache)(nil).OrderBookRowsLength))
}

// PurgeOrderBookRows mocks base method
func (m *MockCache) PurgeOrderBookRows() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PurgeOrderBookRows")
}

// PurgeOrderBookRows indicates an expected call of PurgeOrderBookRows
func (mr *MockCacheMockRecorder) PurgeOrderBookRows() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PurgeOrderBookRows", reflect.TypeOf((*MockCache)(nil).PurgeOrderBookRows))
}
