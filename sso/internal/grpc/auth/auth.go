package auth

import (
	"context"

	validator "github.com/go-playground/validator/v10"
	ssov1 "gitlab.com/kluster1/collection/backend/sso/sso_proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginRequest struct {
	Login    string `validate:"required,login"`
	Password string `validate:"required,password"`
	AppId    int32  `validate:"omitempty"`
}

type RegisterRequest struct {
	Login    string `validate:"required,login"`
	Password string `validate:"required,password"`
}

type Auth interface {
	Login(ctx context.Context, login string, password string, appId int64) (string, error)
	Register(ctx context.Context, login string, password string) (int64, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	body := LoginRequest{req.GetLogin(), req.GetPassword(), req.GetAppId()}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}

	token, err := s.auth.Login(ctx, body.Login, body.Password, int64(body.AppId)) // todo
	if err != nil {
		return nil, status.Error(codes.Internal, "login error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	body := RegisterRequest{req.GetLogin(), req.GetPassword()}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return nil, status.Error(codes.InvalidArgument, "validation error")
	}

	uid, err := s.auth.Register(ctx, body.Login, body.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "register error")
	}

	return &ssov1.RegisterResponse{
		UserId: uid,
	}, nil
}
