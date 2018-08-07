package service

import (
	"github.com/gwbeacon/sdk/v1"
	"google.golang.org/grpc"

	context "golang.org/x/net/context"
)

type QueryService struct{}

func init() {
	Register(&QueryService{})
}

func (s *QueryService) RegisterService(gs *grpc.Server) {
	v1.RegisterQueryServiceServer(gs, s)
}

func (s *QueryService) ServiceVersion() int32 {
	return int32(v1.SdkVersion_V1)
}

func (s *QueryService) ServiceType() int32 {
	return int32(v1.FeatureType_FeatureTypeQuery)
}

func (s *QueryService) GetFeatureList(c context.Context, q *v1.FeatureQuery) (*v1.FeatureList, error) {
	list := make([]*v1.Feature, 0)
	services := GetAllServices()
	for _, service := range services {
		list = append(list, &v1.Feature{
			Type:    v1.FeatureType(service.ServiceType()),
			Name:    TypeToName(service.ServiceType()),
			Version: service.ServiceVersion(),
		})
	}
	return &v1.FeatureList{
		List: list,
	}, nil
}
