package besuConfig

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	*ethclient.Client
}

func New(ctx *context.Context) (*EthClient, error) {
	client, err := ethclient.DialContext(*ctx, os.Getenv("BESU_URL"))
	if err != nil {
		return nil, err
	}

	return &EthClient{
		client,
	}, nil
}
