package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CollectionPairsResponse struct {
	OrderbookPairs []string `json:"orderbook_pairs"`
}

const (
	orderbookURL = "/collection-pairs?candles=false&orderbooks=true"
)

// GetActiveOrderbookPairs gets a list of all the pairs we want to collect data for
func (cr *ConduitRequests) GetActiveOrderbookPairs(retry int) ([]string, error) {

	if retry == 0 {
		return []string{}, fmt.Errorf("Error getting response from gatekeeper")
	}

	client := http.Client{Timeout: cr.cfg.RequestTimeout}

	resp, err := client.Get(cr.cfg.GatekeeperUrl + orderbookURL)
	if err != nil || resp.StatusCode != 200 {
		cr.logger.Infow("Failed getting orderbook pairs from gatekeeper, retrying", "retries_left", retry)
		return cr.GetActiveOrderbookPairs(retry - 1)
	}
	defer cr.closeResponseBody(resp)

	var result CollectionPairsResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []string{}, err
	}
	return result.OrderbookPairs, nil
}
