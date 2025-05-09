package contexts

import (
	"context"

	"github.com/google/uuid"
)

type requestIDKey struct{}

func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey{}, uuid.New().String())
}

func RequestIDFromContext(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		return "", nil
	}
	return requestID, nil
}
