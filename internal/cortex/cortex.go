package cortex

import (
	"context"
	"log"

	"github.com/google/wire"
	"github.com/volatrade/conduit/internal/models"
	client "github.com/volatrade/cortex/external/conduit"
	"google.golang.org/grpc"
)

var Module = wire.NewSet(
	New,
)

type Cortex interface {
	SendOrderBookRow(ob *models.OrderBookRow)
}

type CortexClient struct {
	Client client.ConduitServiceClient
}

func New() (*CortexClient, error) {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
		return nil, err
	}

	con := client.NewConduitServiceClient(conn)

	return &CortexClient{Client: con}, err
}

func (cc *CortexClient) SendOrderBookRow(ob *models.OrderBookRow) {
	cc.Client.HandleOrderBookRow(context.Background(), &client.OrderBookRowRequest{Data: ob.Pair})
}
