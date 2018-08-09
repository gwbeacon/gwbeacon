package service

import (
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
)

type UserService struct {
}

func init() {
	lib.RegisterService(&UserService{})
}

func (s *UserService) Register(gs *grpc.Server) {
	v1.RegisterUserServiceServer(gs, s)
}

func (s *UserService) Version() int32 {
	return int32(v1.SdkVersion_V1)
}

func (s *UserService) Type() int32 {
	return int32(v1.FeatureType_FeatureTypeUser)
}

func (s *UserService) SignUp(ctx context.Context, account *v1.UserAccount) (*v1.Result, error) {
	return nil, nil
}
func (s *UserService) SignIn(ctx context.Context, account *v1.UserAccount) (*v1.Result, error) {
	return nil, nil
}
func (s *UserService) RegisterClient(ctx context.Context, info *v1.ClientInfo) (*v1.Session, error) {
	return nil, nil
}
func (s *UserService) GetInfo(ctx context.Context, user *v1.User) (*v1.UserInfo, error) {
	return nil, nil
}
func (s *UserService) UpdateInfo(ctx context.Context, info *v1.UserInfo) (*v1.UserInfo, error) {
	return nil, nil
}
func (s *UserService) Logout(ctx context.Context, info *v1.UserInfo) (*v1.Result, error) {
	return nil, nil
}
