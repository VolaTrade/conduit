//go:generate mockgen -package=mocks -destination=../mocks/candlerequest.go github.com/volatrade/conduit/internal/candlerequest CandleRequest

package candlerequest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
)

//Module for wire binding
var Module = wire.NewSet(
	New,
)

//Config ...
type Config struct {
	BinanceURL string
}

type CandleRequest interface {
}

type Requests struct {
	mux        *sync.Mutex
	logger     *logger.Logger
	ctx        context.Context
	binanceUrl url.URL
}

func New(ctx context.Context, logger *logger.Logger) *Requests {
	return &Requests{ctx: ctx, logger: logger}
}

func (ccr *CandleRequest) getRecentCandle(symbol string) (*models.Kline, error) {

	resp, err := http.Get(
		fmt.Sprintf("%s?symbol=%s&interval=1m&limit=2",
			ccr.cfg.BinanceApiUrl, symbol),
	)
	if err != nil {
		return nil, err
	}
	defer ccr.closeResponseBody(resp)
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
