package cortex

import (
	"context"
	"fmt"

	"github.com/google/wire"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	conduitpb "github.com/volatrade/protobufs/cortex/conduit"
)

var Module = wire.NewSet(
	New,
)

type (
	Cortex interface {
		SendOrderBookRows(obs []string) error
	}
	Config struct {
		Port int
		Host string
	}
	CortexClient struct {
		//client conduitpb.ConduitServiceClient
		//conn   *grpc.ClientConn
		config *Config
		kstats stats.Stats
		logger *logger.Logger
	}
)

func New(cfg *Config, kstats stats.Stats, logger *logger.Logger) (*CortexClient, error) {

	return &CortexClient{config: cfg, kstats: kstats, logger: logger}, nil
}

func (cc *CortexClient) GetCortexUrlString() string {
	
}

func (cc *CortexClient) SendOrderBookRows(obRows []string) error {

	res, err := cc.client.HandleOrderBookRow(context.Background(), &conduitpb.OrderBookRowRequest{Data: obRows})
	if err != nil {
		cc.kstats.Increment("cortex.errors", 1)
		return fmt.Errorf("response: %+v, error: %s", res, err)
	}
	cc.kstats.Increment("cortex_requests", 1.0)
	return nil
}
