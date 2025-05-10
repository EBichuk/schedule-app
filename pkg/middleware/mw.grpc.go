package middleware

import (
	"context"
	"log/slog"
	"schedule-app/pkg/contexts"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	ctx = contexts.WithRequestID(ctx)

	md, _ := metadata.FromIncomingContext(ctx)
	address, _ := peer.FromContext(ctx)

	slog.InfoContext(ctx, "grpc-request",
		slog.String("remote_addr", address.Addr.String()),
		slog.String("user_agent", md["user-agent"][0]),
	)
	t1 := time.Now()

	m, err := handler(ctx, req)

	slog.InfoContext(ctx, "grpc request completed",
		slog.String("duration", time.Since(t1).String()),
	)

	return m, err
}
