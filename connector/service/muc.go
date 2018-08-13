package service

import (
	"github.com/gwbeacon/sdk/v1"

	context "golang.org/x/net/context"
)

type MucService struct{}

func (muc *MucService) List(ctx context.Context, page *v1.Page) (*v1.MUCList, error) {
	return nil, nil
}
func (muc *MucService) Create(ctx context.Context, info *v1.MUCInfo) (*v1.MUCInfo, error) {
	return nil, nil
}
func (muc *MucService) Destroy(ctx context.Context, info *v1.MUCInfo) (*v1.Result, error) {
	return nil, nil
}
func (muc *MucService) UpdateInfo(ctx context.Context, info *v1.MUCInfo) (*v1.MUCInfo, error) {
	return nil, nil
}
func (muc *MucService) Join(ctx context.Context, info *v1.MUCInfo) (*v1.Result, error) {
	return nil, nil
}
func (muc *MucService) Leave(ctx context.Context, info *v1.UserInfo) (*v1.Result, error) {
	return nil, nil
}
func (muc *MucService) AddMember(ctx context.Context, info *v1.UserInfo) (*v1.Result, error) {
	return nil, nil
}
func (muc *MucService) RemoveMember(ctx context.Context, info *v1.UserInfo) (*v1.Result, error) {
	return nil, nil
}
