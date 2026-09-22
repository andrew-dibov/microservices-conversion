package clients

import (
	"context"
	"microservices-conversion/internal/modules"
	"microservices-conversion/pkg/api/currency"

	"google.golang.org/grpc/metadata"
)

func (client *CurrencyClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

/* --- --- --- */

func (client *CurrencyClient) Rate(ctx context.Context, fromCurrency string, toCurrency string) (*currency.RateResponse, error) {
	if reqID := modules.GetID(ctx); reqID != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", reqID)
	}

	/* --- --- --- */

	ctx, cancel := context.WithTimeout(ctx, client.conf.CurrencyService.Timeout)
	defer cancel()

	return client.grpc.Rate(ctx, &currency.RateRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	})
}

func (client *CurrencyClient) Rates(ctx context.Context, baseCurrency string) (*currency.RatesResponse, error) {
	if reqID := modules.GetID(ctx); reqID != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", reqID)
	}

	/* --- --- --- */

	ctx, cancel := context.WithTimeout(ctx, client.conf.CurrencyService.Timeout)
	defer cancel()

	return client.grpc.Rates(ctx, &currency.RatesRequest{
		BaseCurrency: baseCurrency,
	})
}
