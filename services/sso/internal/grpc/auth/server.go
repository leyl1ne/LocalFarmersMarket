package authgrpc

import (
	"context"
	"errors"

	ssov1 "github.com/leyl1ne/LocalFarmersMarket/pkg/protobuf/go/sso"
	"github.com/leyl1ne/LocalFarmersMarket/services/sso/internal/domain/models"
	"github.com/leyl1ne/LocalFarmersMarket/services/sso/internal/services/auth"
	"github.com/leyl1ne/LocalFarmersMarket/services/sso/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		username string,
		password string,
		phone string,
		role models.UserRole,
		farmName string,
		address models.Address,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	if err := validateRegisterRequest(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var farmName string
	if in.FarmName != nil {
		farmName = *in.FarmName
	}

	var address models.Address
	if in.Address != nil {
		address = models.Address{
			City:     in.Address.GetCity(),
			Street:   in.Address.GetStreet(),
			Building: in.Address.GetBuilding(),
		}
	}

	uid, err := s.auth.RegisterNewUser(ctx,
		in.GetEmail(),
		in.GetUsername(),
		in.GetPassword(),
		in.GetPhone(),
		models.UserRole(in.GetRole().String()),
		farmName,
		address)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: uid}, nil
}

func validateRegisterRequest(in *ssov1.RegisterRequest) error {
	if in.Email == "" {
		return errors.New("email is required")
	}

	if in.Password == "" {
		return errors.New("password is required")
	}

	if in.Username == "" {
		return errors.New("username is required")
	}

	if in.Phone == "" {
		return errors.New("phone is required")
	}

	role := in.GetRole()
	if role == ssov1.UserRole_FARMER {
		if in.FarmName == nil {
			return errors.New("farm_name is required for farmers")
		}

		if in.Address == nil {
			return errors.New("address is required for farmers")
		}
	}

	if role < ssov1.UserRole_CUSTOMER || role > ssov1.UserRole_MODERATOR {
		return errors.New("invalid user role")
	}

	return nil
}
