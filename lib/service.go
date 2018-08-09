package lib

import (
	"log"

	"github.com/gwbeacon/sdk/v1"
	"google.golang.org/grpc"
)

const (
	QueryServiceName   = "query"
	MessageServiceName = "message"
	UserServiceName    = "user"
	RosterServiceName  = "roster"
	MUCServiceName     = "muc"
)

var typeNameMap = map[int32]string{
	int32(v1.FeatureType_FeatureTypeQuery):   QueryServiceName,
	int32(v1.FeatureType_FeatureTypeMessage): MessageServiceName,
	int32(v1.FeatureType_FeatureTypeUser):    UserServiceName,
	int32(v1.FeatureType_FeatureTypeRoster):  RosterServiceName,
	int32(v1.FeatureType_FeatureTypeMUC):     MUCServiceName,
}

func GetServiceInfo(s Service) *ServiceInfo {
	return &ServiceInfo{
		Name:    TypeToName(s.Type()),
		Version: s.Version(),
	}
}

type Service interface {
	Register(gs *grpc.Server)
	Version() int32
	Type() int32
}

var allServices = make([]Service, 0)

func RegisterService(s Service) {
	allServices = append(allServices, s)
}

func LoadAllService(gs *grpc.Server) {
	for _, s := range allServices {
		log.Println("load service ", TypeToName(s.Type()), s.Version())
		s.Register(gs)
	}
}

func GetAllServices() []Service {
	return allServices
}

func GetAllServiceInfo() []*ServiceInfo {
	res := make([]*ServiceInfo, 0)
	for _, s := range GetAllServices() {
		res = append(res, GetServiceInfo(s))
	}
	return res
}

func TypeToName(tp int32) string {
	return typeNameMap[tp]
}
