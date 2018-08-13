package service

import (
	"errors"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"time"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/gwbeacon/lib/rpc"
	"github.com/gwbeacon/sdk/v1"
)

type UserService struct {
	service userService
}

type userService struct {
}

func init() {
	lib.RegisterService(&UserService{})
}

func (s *UserService) Register(gs *grpc.Server) {
	v1.RegisterUserServiceServer(gs, &s.service)
}

func (s *UserService) GetInfo() *rpc.ServiceInfo {
	return &rpc.ServiceInfo{
		Version: int32(v1.SdkVersion_V1),
		Name:    lib.FeatureUserService,
	}
}

func (s *userService) SignUp(ctx context.Context, account *v1.UserAccount) (*v1.Result, error) {
	return nil, nil
}

func (s *userService) SignIn(ctx context.Context, account *v1.UserAccount) (*v1.UserInfo, error) {
	user := &v1.User{
		Domain: account.Domain,
		Name:   account.Name,
	}
	session, ok := ctx.Value(lib.ContextSessionKey).(*rpc.Session)
	if !ok {
		return nil, errors.New("no session")
	}
	connector, ok := ctx.Value(lib.ContextServerKey).(lib.Connector)
	if !ok {
		return nil, errors.New("no connector")
	}
	info, err := s.GetInfo(ctx, user)
	if err != nil {
		return nil, err
	}
	session.User = &rpc.UserInfo{
		Domain:    user.Domain,
		Name:      user.Name,
		NickName:  user.NickName,
		Device:    account.Device,
		LoginTime: uint64(time.Now().Unix()),
	}
	err = connector.Update(session)
	return info, err
}
func (s *userService) RegisterClient(ctx context.Context, info *v1.ClientInfo) (*v1.Session, error) {
	return nil, nil
}
func (s *userService) GetInfo(ctx context.Context, user *v1.User) (*v1.UserInfo, error) {
	return nil, nil
}
func (s *userService) UpdateInfo(ctx context.Context, info *v1.UserInfo) (*v1.UserInfo, error) {
	return nil, nil
}
func (s *userService) Logout(ctx context.Context, info *v1.UserInfo) (*v1.Result, error) {
	return nil, nil
}
