package clients

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (client *RedisClient) Close() error {
	return client.client.Close()
}

/* --- --- --- */

func (client *RedisClient) Get(ctx context.Context, fromCurrency string, toCurrency string) (float64, bool, error) {
	key := fmt.Sprintf("rate:%s:%s",
		strings.ToUpper(strings.TrimSpace(fromCurrency)),
		strings.ToUpper(strings.TrimSpace(toCurrency)))

	rate, err := client.client.Get(ctx, key).Float64()

	if err == redis.Nil {
		return 0, false, nil
	}

	if err != nil {
		return 0, false, fmt.Errorf("RedisClient failed to get rate : %w", err)
	}

	return rate, true, nil
}

func (client *RedisClient) Set(ctx context.Context, fromCurrency string, toCurrency string, rate float64) error {
	key := fmt.Sprintf("rate:%s:%s",
		strings.ToUpper(strings.TrimSpace(fromCurrency)),
		strings.ToUpper(strings.TrimSpace(toCurrency)))

	return client.client.Set(ctx, key, rate, client.ttl).Err()
}
