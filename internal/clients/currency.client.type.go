package clients

import (
	"microservices-conversion/internal/configs"
	"microservices-conversion/pkg/api/currency"

	"google.golang.org/grpc"
)

type CurrencyClient struct {
	grpc currency.CurrencyClient
	conn *grpc.ClientConn
	conf *configs.AppConfig
}
