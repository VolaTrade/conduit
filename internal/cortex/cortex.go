package cortex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
	conduitpb "github.com/volatrade/protobufs/conduit"
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
		Host string
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

	obRow, err := models.UnmarshalDBOrderBookRow(ob)

	rawOb, err := json.Marshal(obRow)
	if err != nil {
		return err
	}
	res, err := cc.client.HandleOrderBookRow(context.Background(),
		&conduitpb.OrderBookRowRequest{Data: rawOb})
	if err != nil {
		println("Are we erroring everytime?")
		return fmt.Errorf("response: %+v, error: %s", res, err)
	}
	cc.kstats.Increment(".cortex_requests", 1.0)
	stringy := fmt.Sprintf("Response from server: %+v", res)
	cc.logger.Infow(stringy)
	return nil
}
