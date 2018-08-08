package service

import (
	"log"
	"net"

	"github.com/gwbeacon/gwbeacon/lib"
	"github.com/gwbeacon/sdk/v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type QueryService struct{}

func init() {
	lib.RegisterService(&QueryService{})
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

func (s *QueryService) GetFeatureList(ctx context.Context, q *v1.FeatureQuery) (*v1.FeatureList, error) {
	p, _ := peer.FromContext(ctx)
	if taddr, ok := p.Addr.(*net.TCPAddr); ok {
		log.Println(taddr.IP.String(), taddr.Port)
	}
	list := make([]*v1.Feature, 0)
	services := lib.GetAllServices()
	for _, service := range services {
		list = append(list, &v1.Feature{
			Type:    v1.FeatureType(service.ServiceType()),
			Name:    lib.TypeToName(service.ServiceType()),
			Version: service.ServiceVersion(),
		})
	}
	return &v1.FeatureList{
		List: list,
	}, nil
}
