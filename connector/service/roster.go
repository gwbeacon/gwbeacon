package service

import (
	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RosterService struct {
}

func init() {
	lib.RegisterService(&RosterService{})
}

func (s *RosterService) Register(gs *grpc.Server) {
	v1.RegisterRosterServiceServer(gs, s)
}

func (s *RosterService) Version() int32 {
	return int32(v1.SdkVersion_V1)
}

func (s *RosterService) Type() int32 {
	return int32(v1.FeatureType_FeatureTypeRoster)
}

func (s *RosterService) List(ctx context.Context, page *v1.Page) (*v1.RosterList, error) {
	return nil, nil
}
func (s *RosterService) Add(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
func (s *RosterService) Remove(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
func (s *RosterService) AddBlock(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
func (s *RosterService) RemoveBlock(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
func (s *RosterService) AddBlack(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
func (s *RosterService) RemoveBlack(context.Context, *v1.User) (*v1.Result, error) {
	return nil, nil
}
