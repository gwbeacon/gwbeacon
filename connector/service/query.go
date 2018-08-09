package service

import (
	"log"

	"github.com/gwbeacon/gwbeacon/lib"
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

func (s *QueryService) Version() int32 {
	return int32(v1.SdkVersion_V1)
}

func (s *QueryService) Type() int32 {
	return int32(v1.FeatureType_FeatureTypeQuery)
}

func (s *QueryService) GetFeatureList(ctx context.Context, q *v1.FeatureQuery) (*v1.FeatureList, error) {
	sess, ok := ctx.Value("session").(*lib.Session)
	if ok {
		log.Println(sess.ID)
	}
	list := make([]*v1.Feature, 0)
	services := lib.GetAllServices()
	for _, service := range services {
		list = append(list, &v1.Feature{
			Type:    v1.FeatureType(service.Type()),
			Name:    lib.TypeToName(service.Type()),
			Version: service.Version(),
		})
	}
	return &v1.FeatureList{
		List: list,
	}, nil
}
