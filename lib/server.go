package lib

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"sync"

	"errors"
	"net"

	"time"

	"github.com/gwbeacon/gwbeacon/lib/rpc"
)

type ServerType string

const (
	RegisterServer  ServerType = "register"
	ConnectorServer            = "connector"
	UserServer                 = "user"
	ChatServer                 = "chat"
	RosterServer               = "roster"
	MUCServer                  = "muc"
)

type Server interface {
	GetInfo() *rpc.ServerInfo
	Register(regAddr string, onIDChange func(uint32)) uint32
	Serve(opt ...grpc.ServerOption) error
	Stop()
}

type server struct {
	sync.RWMutex
	info        rpc.ServerInfo
	regAddr     string
	serversInfo map[string]map[uint32]*rpc.ServerInfo
	services    []Service
	grpcServer  *grpc.Server
	onIDChange  func(uint32)
	registered  bool
}

type ServerInfo struct {
	Type string
	Port int32
}

func NewServer(info ServerInfo, services []Service) Server {
	s := &server{
		info: rpc.ServerInfo{
			Type:     info.Type,
			Port:     info.Port,
			Services: make([]*rpc.ServiceInfo, 0),
		},
		serversInfo: make(map[string]map[uint32]*rpc.ServerInfo),
		services:    services,
	}
	for _, service := range services {
		s.info.Services = append(s.info.Services, service.GetInfo())
	}
	return s
}

func (s *server) GetInfo() *rpc.ServerInfo {
	s.Lock()
	defer s.Unlock()
	return &s.info
}

func (s *server) Register(regAddr string, onIDChange func(uint32)) uint32 {
	s.Lock()
	if s.registered {
		s.Unlock()
		return s.info.ID
	}
	s.regAddr = regAddr
	idNotify := make(chan uint32, 1)
	s.onIDChange = func(id uint32) {
		onIDChange(id)
		idNotify <- id
		close(idNotify)
	}
	s.Unlock()
	go func() {
		for {
			err := s.register()
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
	id := <-idNotify
	s.Lock()
	s.onIDChange = onIDChange
	s.Unlock()
	return id
}

func (s *server) Serve(opt ...grpc.ServerOption) error {
	addr := fmt.Sprintf(":%d", s.info.Port)
	fmt.Println(addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()
	server := grpc.NewServer(opt...)
	if server == nil {
		return errors.New("create grpc server failed")
	}
	s.Lock()
	for _, service := range s.services {
		service.Register(server)
	}
	log.Println(server.GetServiceInfo())
	s.grpcServer = server
	s.Unlock()
	log.Println("start server", s.info)
	err = server.Serve(lis)
	return err
}

func (s *server) Stop() {
	s.Lock()
	defer s.Unlock()
	s.grpcServer.GracefulStop()
}

func (s *server) register() error {
	s.RLock()
	regAddr := s.regAddr
	serverInfo := &s.info
	s.RUnlock()
	conn, err := grpc.DialContext(context.Background(), regAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	fmt.Println("begin to register server", s.info)
	client := rpc.NewClusterClient(conn)
	stream, err := client.Register(context.Background())
	if err != nil {
		conn.Close()
		return err
	}
	err = stream.Send(serverInfo)
	if err != nil {
		conn.Close()
		return err
	}
	for {
		result, err := stream.Recv()
		if err != nil {
			conn.Close()
			return err
		}
		fmt.Println("get server status from register", result)
		s.Lock()
		if result.ID > 0 {
			s.info.ID = result.ID
			s.onIDChange(result.ID)
		}
		for tp, idMap := range result.DownServers {
			if idMap == nil || len(idMap.Servers) == 0 {
				continue
			}
			for _, info := range idMap.Servers {
				if _, ok := s.serversInfo[tp]; ok {
					if _, ok := s.serversInfo[tp][info.ID]; ok {
						delete(s.serversInfo[tp], info.ID)
					}
				}
			}
		}

		for tp, idMap := range result.NewServers {
			s.serversInfo[tp] = idMap.Servers
		}
		s.Unlock()
	}
	return nil
}
