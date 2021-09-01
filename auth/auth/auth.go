package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	OpenIDResolver OpenIDResolver
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

type OpenIDResolver interface {
	Resolve(code string) (string, error)
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

	return &authpb.LoginResponse{
		AccessToken: "token for " + accountID,
		ExpiresIn:   7200,
	}, nil
}
