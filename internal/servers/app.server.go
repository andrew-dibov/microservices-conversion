package servers

import (
	"context"
	"microservices-conversion/internal/clients"
	"microservices-conversion/internal/loggers"
	"microservices-conversion/internal/modules"
	"microservices-conversion/pkg/api/conversion"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAppServer(currencyClient *clients.CurrencyClient, redisClient *clients.RedisClient, appLogger *loggers.AppLogger) *AppServer {
	return &AppServer{
		currencyClient: currencyClient,
		redisClient:    redisClient,
		appLogger:      appLogger,
	}
}

/* --- --- --- */

func (server *AppServer) Convert(ctx context.Context, request *conversion.ConvertRequest) (*conversion.ConvertResponse, error) {
	if request.FromCurrency == "" || request.ToCurrency == "" {
		return nil, status.Error(codes.InvalidArgument, "FromCurrency or ToCurrency is empty")
	}

	if len(request.FromCurrency) != 3 || len(request.ToCurrency) != 3 {
		return nil, status.Error(codes.InvalidArgument, "FromCurrency or ToCurrency is not 3 chars")
	}

	for _, char := range request.FromCurrency {
		if char < 'A' || char > 'Z' {
			return nil, status.Error(codes.InvalidArgument, "FromCurrency has invalid chars")
		}
	}

	for _, char := range request.ToCurrency {
		if char < 'A' || char > 'Z' {
			return nil, status.Error(codes.InvalidArgument, "ToCurrency has invalid chars")
		}
	}

	if request.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Amount must be positive")
	}

	if request.Amount > 1e12 {
		return nil, status.Error(codes.InvalidArgument, "Amount is too large")
	}

	/* --- --- --- */

	rate, err := server.rate(ctx, request.FromCurrency, request.ToCurrency)
	if err != nil {
		return nil, err
	}

	/* --- --- --- */

	server.appLogger.Info("Convert completed",
		"from_currency", request.FromCurrency,
		"to_currency", request.ToCurrency,
		"amount", request.Amount,
		"result", request.Amount*rate,
		"rate", rate,
		"request_id", modules.GetID(ctx))

	response := &conversion.ConvertResponse{
		FromCurrency: request.FromCurrency,
		ToCurrency:   request.ToCurrency,
		Amount:       request.Amount,
		Result:       request.Amount * rate,
		Rate:         rate,
	}

	/* --- --- --- */

	/* push event to kafka for history service */

	/* --- --- --- */

	return response, nil
}

/* --- --- --- */

func (server *AppServer) rate(ctx context.Context, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	/* --- --- --- */

	if rate, ok, err := server.redisClient.Get(ctx, fromCurrency, toCurrency); err == nil {
		if ok {
			return rate, nil
		}
	} else {
		server.appLogger.Error("Convert failed to get cache", "error", err)
	}

	/* --- --- --- */

	rateResponse, err := server.currencyClient.Rate(ctx, fromCurrency, toCurrency)
	if err != nil {
		server.appLogger.Error("AppServer get rate failed", "error", err,
			"from_currency", fromCurrency,
			"to_currency", toCurrency,
			"request_id", modules.GetID(ctx))
		return 0, err
	}

	if rateResponse.Rate <= 0 {
		server.appLogger.Error("AppServer received invalid rate", "rate", rateResponse.Rate,
			"request_id", modules.GetID(ctx))
		return 0, status.Error(codes.Internal, "AppServer received invalid rate")
	}

	/* --- --- --- */

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // move to config
		defer cancel()

		if err := server.redisClient.Set(ctx, fromCurrency, toCurrency, rateResponse.Rate); err != nil {
			server.appLogger.Error("Convert failed to set cache", "error", err)
		}
	}()

	/* --- --- --- */

	return rateResponse.Rate, nil
}
