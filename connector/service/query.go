package service

import (
	"log"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/gwbeacon/lib/rpc"
	"github.com/gwbeacon/sdk/v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type QueryService struct{}

func init() {
	lib.RegisterService(&QueryService{})
}

func (s *QueryService) Register(gs *grpc.Server) {
	v1.RegisterQueryServiceServer(gs, s)
}

func (s *QueryService) GetInfo() *rpc.ServiceInfo {
	return &rpc.ServiceInfo{
		Version: int32(v1.SdkVersion_V1),
		Name:    lib.FeatureQueryService,
	}
}

func (s *QueryService) GetFeatureList(ctx context.Context, q *v1.FeatureQuery) (*v1.FeatureList, error) {
	sess, ok := ctx.Value("session").(*rpc.Session)
	if ok {
		log.Println(sess.ID)
	}
	list := make([]*v1.Feature, 0)
	services := lib.GetAllServices()
	for _, service := range services {
		info := service.GetInfo()
		list = append(list, &v1.Feature{
			Name:    info.Name,
			Version: info.Version,
		})
	}
	return &v1.FeatureList{
		List: list,
	}, nil
}
