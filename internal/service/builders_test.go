package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOrderBookChannels(t *testing.T) {

	ts := createTestSuite(t)
	ts.service.BuildOrderBookChannels(3)
	c := ts.service.GetOrderBookChannel(2)
	assert.True(t, c != nil)
}

func TestBuildTransactionChannels(t *testing.T) {
	ts := createTestSuite(t)
	ts.service.BuildTransactionChannels(3)
	c := ts.service.GetTransactionChannel(2)
	assert.True(t, c != nil)
}

func TestSpawnSocketRoutines(t *testing.T) {
	ts := createTestSuite(t)

	sockets := ts.service.SpawnSocketRoutines(3)
	assert.True(t, sockets != nil)
}
