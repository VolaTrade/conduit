package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/volatrade/conduit/internal/models"
)

const (
	pedersonURL = "/v1/update/pederson"
)

//PostOrderbookRow sends a POST request to Cortex to update it with the most recent orderbook data
func (cr *ConduitRequests) PostOrderbookRowToCortex(orderbookRow *models.OrderBookRow) error {

	postUrl := fmt.Sprintf("%s:%d%s", cr.cfg.CortexUrl, cr.cfg.CortexPort, pedersonURL)

	cr.logger.Infow("Sending request to Cortex", "url", postUrl)
	data, err := json.Marshal(orderbookRow)
	if err != nil {
		return err
	}

	resp, err := http.Post(postUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		cr.statz.Increment("cortex.errors", 1)
		return fmt.Errorf("response error: %s", err.Error())
	}

	if err := resp.Body.Close(); err != nil {
		cr.logger.Errorw("Error closing response: ", "error", err)
	}

	cr.logger.Infow(fmt.Sprintf("Cortex request success, response: %s", resp.Header))
	cr.statz.Increment("cortex_requests", 1.0)

	return nil
}
