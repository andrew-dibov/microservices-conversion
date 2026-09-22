package modules

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetStringEnv(key string, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func GetIntegerEnv(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil {
			return val
		}
	}
	return def
}

func GetBooleanEnv(key string, def bool) bool {
	if str := os.Getenv(key); str != "" {
		if val, err := strconv.ParseBool(str); err == nil {
			return val
		}
	}
	return def
}

func GetDurationEnv(key string, def time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		if val, err := time.ParseDuration(str); err == nil {
			return val
		}
	}
	return def
}

func GetStringSetEnv(key string, def map[string]bool) map[string]bool {
	str := os.Getenv(key)
	keys := make(map[string]bool)

	if str != "" {
		for _, k := range strings.Split(str, ",") {
			keys[strings.TrimSpace(k)] = true
		}
	} else {
		keys = def
	}

	return keys
}
