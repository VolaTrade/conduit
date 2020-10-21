package client_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/volatrade/candles/internal/client"
	"github.com/volatrade/candles/internal/mocks"
)

type testSuite struct {
	mockController *gomock.Controller
	mockStats      *mocks.MockStats
	client         *client.ApiClient
}

func createTestSuite(t *testing.T) testSuite {
	var err error

	mockController := gomock.NewController(t)
	mockStats := mocks.NewMockStats(mockController)

	testClient := client.New(mockStats)
	assert.NoError(t, err)

	return testSuite{
		mockController: mockController,
		mockStats:      mockStats,
		client:         testClient,
	}
}

func TestMain(m *testing.M) {

	retCode := m.Run()

	os.Exit(retCode)
}
