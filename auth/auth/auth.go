package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	OpenIDResolver OpenIDResolver
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

func (s *Service) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	//s.Logger.Info("received code", zap.String("code", req.Code))
	openID, err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid: %v", err)
	}

	accountID, err := s.Mongo.ResolveAccountID(ctx, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	tkn, err := s.TokenGenerator.GenerateToken(accountID, s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "")
	}
	return &authpb.LoginResponse{
		AccessToken: "token for " + tkn,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
