package grpc

import (
	"context"

	"github.com/018bf/example/pkg/log"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type RequestIDMiddleware struct {
}

func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

func (m *RequestIDMiddleware) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	newCtx := context.WithValue(ctx, log.RequestIDKey, uuid.New().String())
	return handler(newCtx, req)
}
