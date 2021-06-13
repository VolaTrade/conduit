//go:generate mockgen -package=mocks -destination=../mocks/candlerequest.go github.com/volatrade/conduit/internal/candlerequest CandleRequest

package candlerequest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
)

// //Module for wire binding
// var Module = wire.NewSet(
// 	New,
// )

//Config ...
type Config struct {
	BinanceUrl string
}

// type CandleRequest interface {
// 	getRecentCandle(string) (*models.Kline, error)
// 	closeResponseBody() *http.Response
// }

type CandleRequest struct {
	mux    *sync.Mutex
	logger *logger.Logger
	ctx    context.Context
	cfg    *Config
}

func New(ctx context.Context, logger *logger.Logger) *CandleRequest {
	return &CandleRequest{ctx: ctx, logger: logger}
}

func (cr *CandleRequest) getRecentCandle(symbol string) (*models.Kline, error) {

	resp, err := http.Get(
		fmt.Sprintf("%s?symbol=%s&interval=1m&limit=2",
			cr.cfg.BinanceUrl, symbol),
	)
	if err != nil {
		return nil, err
	}
	defer cr.closeResponseBody(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got unexpected status code : % d , expected 200", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	kline, err := models.UnmarshallBytesToLatestKline(body, symbol)
	if err != nil {
		return nil, err
	}

	return kline, nil
}

func (cr *CandleRequest) closeResponseBody(resp *http.Response) {

	if err := resp.Body.Close(); err != nil {
		cr.logger.Errorw("Error closing response")
	}
}
