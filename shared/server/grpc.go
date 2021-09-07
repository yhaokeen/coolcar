package server

import (
	"coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	RegisterFunc      func(server *grpc.Server)
	Logger            *zap.Logger
}

func RunGRPCServer(cfg *GRPCConfig) error {
	nameField := zap.String("name", cfg.Name)
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		cfg.Logger.Fatal("cannot listen", nameField, zap.Error(err))
	}
	var opts []grpc.ServerOption
	if cfg.AuthPublicKeyFile != "" {
		in, err := auth.Interceptor(cfg.AuthPublicKeyFile)
		if err != nil {
			cfg.Logger.Fatal("cannot create auth interceptor")
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}
	s := grpc.NewServer(opts...)

	cfg.RegisterFunc(s)
	cfg.Logger.Info("server started", nameField, zap.String("addr", cfg.Addr))
	return s.Serve(lis)
}
