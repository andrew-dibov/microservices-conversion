package servers

import (
	"microservices-conversion/internal/clients"
	"microservices-conversion/internal/loggers"
	"microservices-conversion/pkg/api/conversion"
)

type AppServer struct {
	conversion.UnimplementedConversionServer
	currencyClient *clients.CurrencyClient
	appLogger      *loggers.AppLogger
}
