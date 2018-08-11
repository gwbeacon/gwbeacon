package lib

import (
	"log"

	"github.com/gwbeacon/gwbeacon/lib/rpc"
	"google.golang.org/grpc"
)

const (
	FeatureConnector      = "connector"
	FeatureQueryService   = "query"
	FeatureMessageService = "message"
	FeatureUserService    = "user"
	FeatureRosterService  = "roster"
	FeatureMUCService     = "muc"
)

type Service interface {
	Register(gs *grpc.Server)
	GetInfo() *rpc.ServiceInfo
}

var allServices = make([]Service, 0)

func RegisterService(s Service) {
	allServices = append(allServices, s)
}

func LoadAllService(gs *grpc.Server) {
	for _, s := range allServices {
		log.Println("load service ", s.GetInfo())
		s.Register(gs)
	}
}

func GetAllServices() []Service {
	return allServices
}

func GetAllServiceInfo() []*rpc.ServiceInfo {
	res := make([]*rpc.ServiceInfo, 0)
	for _, s := range GetAllServices() {
		res = append(res, s.GetInfo())
	}
	return res
}
