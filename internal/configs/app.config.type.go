package configs

import "time"

type AppConfig struct {
	App App

	CurrencyService CurrencyService

	RedisClient RedisClient
}

/* --- --- --- */

type App struct {
	Name string

	Prod bool
	Port string

	KeepaliveTime    time.Duration
	KeepaliveTimeout time.Duration

	ShutdownTimeout time.Duration
}

/* --- --- --- */

type CurrencyService struct {
	Address string
	Timeout time.Duration
}

/* --- --- --- */

type RedisClient struct {
	Address string

	Database int
	Password string
	TTL      time.Duration

	Retries RedisRetries
}

type RedisRetries struct {
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
}
