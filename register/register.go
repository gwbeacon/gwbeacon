package main

import (
	"flag"
	"log"
	"net"
	"sync"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/gwbeacon/gwbeacon/lib"
)

type server struct {
	sync.Mutex
	ids     map[string]uint32
	servers map[interface{}]*serverInfo
}

type serverInfo struct {
	info   *lib.ServerInfo
	stream lib.Cluster_RegisterServer
}

func NewServer() *server {
	return &server{
		ids:     make(map[string]uint32),
		servers: make(map[interface{}]*serverInfo),
	}
}

func (s *server) GetAllServers() *lib.RegisterReturn {
	s.Lock()
	defer s.Unlock()
	res := &lib.RegisterReturn{
		NewServers: make([]*lib.ServerInfo, 0),
	}
	for _, info := range s.servers {
		res.NewServers = append(res.NewServers, info.info)
	}
	return res
}

func (s *server) Register(info *lib.ServerInfo, stream lib.Cluster_RegisterServer) error {
	s.Lock()
	ctx := stream.Context()
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		log.Println("1")
		return nil
	}
	if addr, ok := p.Addr.(*net.TCPAddr); ok {
		info.IP = addr.IP.String()
		info.Port = int32(addr.Port)
	}
	if _, ok := s.servers[ctx]; ok {
		log.Println("2")
		return nil
	}
	if _, ok := s.ids[info.Type]; !ok {
		s.ids[info.Type] = 1
	}
	info.ID = s.ids[info.Type]
	s.ids[info.Type]++
	info1 := &serverInfo{
		info:   info,
		stream: stream,
	}
	log.Println("register server", info)
	s.servers[ctx] = info1
	s.Unlock()
	stream.Send(s.GetAllServers())
	return nil
}

func (s *server) Sync(ctx context.Context, info *lib.ServerInfo) (*lib.RegisterReturn, error) {
	return nil, nil
}

func main() {
	var port = ""
	flag.StringVar(&port, "port", "9999", "-port 9999")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return
	}

	server := grpc.NewServer()
	lib.RegisterClusterServer(server, NewServer())
	err = server.Serve(lis)
	log.Println(err)
}
