package cortex

import (
	"context"
	"fmt"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	conduitpb "github.com/volatrade/cortex/external/conduit"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	"google.golang.org/grpc"
)

var Module = wire.NewSet(
	New,
)

type (
	Cortex interface {
		SendOrderBookRow(ob *models.OrderBookRow) error
	}

	Config struct {
		Port int
	}

	CortexClient struct {
		client conduitpb.ConduitServiceClient
		conn   *grpc.ClientConn
		config *Config
		kstats *stats.Stats
		logger *logger.Logger
	}
)

func New(cfg *Config, kstats *stats.Stats, logger *logger.Logger) (*CortexClient, func(), error) {

	log.Println("creating client connection to cortex -> port:", cfg.Port)
	conn, err := grpc.Dial(fmt.Sprintf(":%d", cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %s", err)
		return nil, nil, err
	}
	client := conduitpb.NewConduitServiceClient(conn)

	end := func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing client connection to cortex: %v", err)
			}
			log.Println("Successful Shutdown of client connection to cortex")
		}
	}

	return &CortexClient{client: client, conn: conn, config: cfg, kstats: kstats, logger: logger}, end, nil
}

func (cc *CortexClient) SendOrderBookRow(ob *models.OrderBookRow) error {
	res, err := cc.client.HandleOrderBookRow(context.Background(), &conduitpb.OrderBookRowRequest{Data: ob.Pair})
	if err != nil {
		return fmt.Errorf("%+v.%s", res, err)
	}
	cc.kstats.Increment(fmt.Sprintf("%+v", res), 1.0)
	log.Printf("Response from server: %+v", res)
	return nil
}
