package configs

import "time"

type AppConfig struct {
	App App

	CurrencyService CurrencyService
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
