package interceptors

import (
	"context"
	"microservices-conversion/internal/modules"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceInterceptor(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	meta, ok := metadata.FromIncomingContext(ctx)

	if ok {
		if ids := meta.Get("x-request-id"); len(ids) > 0 {
			ctx = context.WithValue(ctx, modules.ID{}, ids[0])
		}
	}

	return handler(ctx, request)
}
