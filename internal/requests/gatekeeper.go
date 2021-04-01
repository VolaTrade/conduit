package requests

import (
	"encoding/json"
	"net/http"
)

const VT_COLLECTION_PAIRS_URL = "https://kjabyesfa1.execute-api.us-west-2.amazonaws.com/prod/collection-pairs?"

type CollectionPairsResponse struct {
	OrderbookPairs []string `json:"orderbook_pairs"`
}

// GetActiveOrderbookPairs gets a list of all the pairs we want to collect data for
func (ac *ConduitRequests) GetActiveOrderbookPairs() ([]string, error) {
	resp, err := http.Get(VT_COLLECTION_PAIRS_URL + "transaction=false&orderbook=true")
	if err != nil {
		return nil, err
	}

	var result CollectionPairsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.OrderbookPairs, nil
}
