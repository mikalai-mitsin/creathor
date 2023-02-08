package rest

import (
	"context"
	"github.com/018bf/example/internal/configs"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"

	"github.com/018bf/example/pkg/log"
)

type Server struct {
	server *gin.Engine
	config *configs.Config
	logger log.Logger
}

func NewServer(
	logger log.Logger,
	config *configs.Config,
) *Server {
	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Grpc-Metadata-Count")
	corsConfig.AddExposeHeaders("Grpc-Metadata-Count")
	server.Use(cors.New(corsConfig))
	return &Server{
		server: server,
		config: config,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	_ = examplepb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterEquipmentServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterPlanServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterDayServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	_ = examplepb.RegisterArchServiceHandlerFromEndpoint(ctx, mux, s.config.BindAddr, opts)
	s.server.Any("/*any", gin.WrapH(mux))
	return http.ListenAndServe(":8001", s.server)
}

func (s *Server) Stop(_ context.Context) error {
	return nil
}
