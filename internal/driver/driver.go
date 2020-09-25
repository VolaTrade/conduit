package driver

import (
	"fmt"
	"time"

	"github.com/google/wire"
	"github.com/volatrade/candles/internal/cache"
	"github.com/volatrade/candles/internal/dynamo"
	"github.com/volatrade/candles/internal/models"
	"github.com/volatrade/candles/internal/service"
)

var Module = wire.NewSet(
	New,
)

type (
	CandlesDriver struct {
		svc *service.CandlesService
	}
)

type Driver interface {
	Run()
}

func New(service *service.CandlesService) *CandlesDriver {
	return &CandlesDriver{svc: service}
}

func (cd *CandlesDriver) Run() {

	if err := cd.svc.Init(); err != nil {
		panic(err)
	}
	//Insert concurrent workload distribution here

	// Test insertion
	candle, err := cache.NewCandle("1234", "1245", "1245", "12455")
	if err != nil {
		panic(err)
	}

	dynamo, err := dynamo.New(&dynamo.Config{TableName: "candles"})
	if err != nil {
		panic(err)
	}

	pairData := cache.InitializePairData()
	pairData.Five = []*models.Candle{candle}

	dynamoItem := &models.DynamoCandleItem{
		PairData:  pairData,
		PairName:  "ETHUSDT",
		Timestamp: time.Now().String(),
	}

	tableStatus, err := dynamo.CreateCandlesTable()
	if err != nil {
		panic(err)
	}

	isHealthy, err := dynamo.IsHealthy()
	if err != nil {
		panic(err)
	}

	fmt.Println("HEALTH CHECK: ", isHealthy)

	for isHealthy == false {
		ms := time.Now().Nanosecond()
		s := time.Now().Second()

		// Every 5 seconds check table status
		if s%5 == 0 && ms == 0 {
			fmt.Printf("Table status: %+v\n", tableStatus)
			isHealthy, err = dynamo.IsHealthy()
			if err != nil {
				panic(err)
			}
		}

		// timeout after 30 seconds
		if s%30 == 0 {
			panic("Table creation timed out, took longer than 30 seconds")
		}

	}

	err = dynamo.AddItem(dynamoItem)
	if err != nil {
		panic(err)
	}

}
