package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/volatrade/conduit/internal/models"
)

func (cr *ConduitRequests) GetRecentCandle(symbol string) (*models.Kline, error) {

	resp, err := http.Get(
		fmt.Sprintf("%s?symbol=%s&interval=1m&limit=2",
			cr.cfg.BinanceApiUrl, symbol),
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
	kline, err := models.UnmarshallBytesToLatestKline(body)
	if err != nil {
		return nil, err
	}

	return kline, nil
}
