package grpc

import (
	"context"
	"github.com/018bf/example/internal/configs"
	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"strings"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ctxKey int

const (
	UserKey ctxKey = iota + 1
)

type AuthMiddleware struct {
	logger          log.Logger
	config          *configs.Config
	authInterceptor interceptors.AuthInterceptor
}

func NewAuthMiddleware(
	authInterceptor interceptors.AuthInterceptor,
	logger log.Logger,
	config *configs.Config,
) *AuthMiddleware {
	return &AuthMiddleware{authInterceptor: authInterceptor, logger: logger, config: config}
}

func (m *AuthMiddleware) Auth(ctx context.Context) (context.Context, error) {
	var token models.Token
	stringToken, err := AuthFromMD(ctx)
	if err != nil {
		return context.WithValue(ctx, UserKey, models.Guest), nil
	}
	token = models.Token(stringToken)
	user, err := m.authInterceptor.Auth(ctx, token)
	if err != nil {
		return nil, decodeError(err)
	}
	newCtx := context.WithValue(ctx, UserKey, user)
	return newCtx, nil
}

func (m *AuthMiddleware) UnaryServerInterceptorAuth(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	newCtx, err := m.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}

func RequestIDUnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx := context.WithValue(ctx, log.RequestIDKey, uuid.New().String())
	return handler(newCtx, req)
}

const (
	headerAuthorize = "authorization"
	expectedScheme  = "bearer"
)

func AuthFromMD(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}
	return splits[1], nil
}
