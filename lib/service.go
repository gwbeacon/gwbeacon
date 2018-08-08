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

type Service interface {
	RegisterService(gs *grpc.Server)
	ServiceVersion() int32
	ServiceType() int32
}

var allServices = make([]Service, 0)

func RegisterService(s Service) {
	allServices = append(allServices, s)
}

func LoadAllService(gs *grpc.Server) {
	for _, s := range allServices {
		log.Println("load service ", TypeToName(s.ServiceType()), s.ServiceVersion())
		s.RegisterService(gs)
	}
}

func GetAllServices() []Service {
	return allServices
}

func TypeToName(tp int32) string {
	return typeNameMap[tp]
}
