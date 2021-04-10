package cortex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

const (
	CORTEX_PEDERSON_URL string = "/v1/pederson"
)

type (
	Cortex interface {
		SendFullCacheUpdate(string) error
	}
	Config struct {
		Port int
		Host string
	}
	CortexConnection struct {
		config *Config
		kstats stats.Stats
		logger *logger.Logger
		url    string
	}
)

func New(cfg *Config, kstats stats.Stats, logger *logger.Logger) (*CortexConnection, error) {

	cortexUpdateUrl := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	innerPath := fmt.Sprintf("%s", CORTEX_PEDERSON_URL)
	pedersonUrl := url.URL{Scheme: "http", Host: cortexUpdateUrl, Path: innerPath}
	return &CortexConnection{url: pedersonUrl.String(), config: cfg, kstats: kstats, logger: logger}, nil
}

// func (cc *CortexConnection) GetCortexUrlString() string {

// 	// log.Printf("%s", updateUrl.String())
// 	return updateUrl.String()
// }

func (cc *CortexConnection) SendFullCacheUpdate(pair string) error {

	var message models.CortexRequest
	message.Pair = pair
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(cc.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		cc.kstats.Increment("cortex.errors", 1)
		return fmt.Errorf("response: %+v, error: %s", resp.Status, err)
	}

	cc.logger.Infow(fmt.Sprintf("Cortex request success, response: %s", resp.Header))
	cc.kstats.Increment("cortex_requests", 1.0)
	return nil
}
