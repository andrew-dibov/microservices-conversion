package configs

import (
	"microservices-conversion/internal/modules"
	"time"
)

func NewAppConfig() AppConfig {
	return AppConfig{
		App: App{
			Name: modules.GetStringEnv("APP_NAME", "microservices-conversion"),

			Prod: modules.GetBooleanEnv("APP_PROD", false),
			Port: modules.GetStringEnv("APP_PORT", "50053"),

			KeepaliveTime:    modules.GetDurationEnv("APP_KEEPALIVE_TIME", 5*time.Second),
			KeepaliveTimeout: modules.GetDurationEnv("APP_KEEPALIVE_TIMEOUT", 5*time.Second),

			ShutdownTimeout: modules.GetDurationEnv("APP_SHUTDOWN_TIMEOUT", 5*time.Second),
		},

		CurrencyService: CurrencyService{
			Address: modules.GetStringEnv("CURRENCY_ADDRESS", "localhost:50052"),
			Timeout: modules.GetDurationEnv("CURRENCY_TIMEOUT", 5*time.Second),
		},
	}
}
