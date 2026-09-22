package modules

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

func GenID() string {
	id := make([]byte, 16)

	if _, err := rand.Read(id); err != nil {
		return "00000000000000000000000000000000"
	}

	return hex.EncodeToString(id)
}

func GetID(ctx context.Context) string {
	if val := ctx.Value(ID{}); val != nil {
		return val.(string)
	}

	return ""
}
