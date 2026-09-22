package clients

import (
	"fmt"
	"microservices-conversion/internal/configs"
	"microservices-conversion/pkg/api/currency"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewCurrencyClient(appConfig *configs.AppConfig) (*CurrencyClient, error) {
	connection, err := grpc.NewClient(appConfig.CurrencyService.Address,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4*1024*1024), grpc.MaxCallSendMsgSize(4*1024*1024)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10 * time.Second, Timeout: 1 * time.Second}),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 2 * time.Second}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
		"methodConfig": [{ "name": [{"service": "currency.Currency"}],
			"retryPolicy": {
				"maxAttempts": 3,
				"maxBackoff": "1s",
				"backoffMultiplier": 2,
				"initialBackoff": "0.1s",
				"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`))

	if err != nil {
		return nil, fmt.Errorf("failed to create CurrencyClient : %w", err)
	}

	return &CurrencyClient{
		grpc: currency.NewCurrencyClient(connection),
		conn: connection,
		conf: appConfig,
	}, nil
}
