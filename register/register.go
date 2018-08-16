package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/gwbeacon/gwbeacon/lib/rpc"
)

type server struct {
	sync.Mutex
	timeBase int32
	ids      map[string]uint32
	servers  map[string]*rpc.ServerInfoIDMap
	streams  map[string]*StreamInfo
}

type StreamInfo struct {
	info       *rpc.ServerInfo
	isDown     bool
	serverUp   chan *rpc.ServerInfo
	serverDown chan *rpc.ServerInfo
}

func NewServer(timeBase int64) *server {
	return &server{
		timeBase: int32(timeBase),
		ids:      make(map[string]uint32),
		servers:  make(map[string]*rpc.ServerInfoIDMap),
		streams:  make(map[string]*StreamInfo),
	}
}

func (s *server) GetAllServers() *rpc.RegisterReturn {
	s.Lock()
	defer s.Unlock()
	return &rpc.RegisterReturn{
		NewServers: s.getAllServers(),
	}
}

func (s *server) getAllServers() map[string]*rpc.ServerInfoIDMap {
	ret := make(map[string]*rpc.ServerInfoIDMap)
	for tp, idMap := range s.servers {
		servers := idMap.Servers
		idMap = new(rpc.ServerInfoIDMap)
		idMap.Servers = make(map[uint32]*rpc.ServerInfo)
		for _, info := range servers {
			addr := fmt.Sprintf("%s:%d", info.IP, info.Port)
			if stream, ok := s.streams[addr]; ok {
				if stream.isDown == false {
					idMap.Servers[info.ID] = info
				}
			}
		}
		if len(idMap.Servers) > 0 {
			ret[tp] = idMap
		}
	}
	return ret
}

func (s *server) onServerUp(stream rpc.Cluster_RegisterServer, info *rpc.ServerInfo) {
	ret := &rpc.RegisterReturn{
		NewServers: make(map[string]*rpc.ServerInfoIDMap),
	}
	ret.NewServers[info.Type] = &rpc.ServerInfoIDMap{
		Servers: make(map[uint32]*rpc.ServerInfo),
	}
	ret.NewServers[info.Type].Servers[info.ID] = info
	err := stream.Send(ret)
	log.Println(err)
}

func (s *server) onServerDown(stream rpc.Cluster_RegisterServer, info *rpc.ServerInfo) {
	ret := &rpc.RegisterReturn{
		DownServers: make(map[string]*rpc.ServerInfoIDMap),
	}
	ret.DownServers[info.Type] = &rpc.ServerInfoIDMap{
		Servers: make(map[uint32]*rpc.ServerInfo),
	}
	ret.DownServers[info.Type].Servers[info.ID] = info
	err := stream.Send(ret)
	log.Println(err)
}

func (s *server) Register(stream rpc.Cluster_RegisterServer) error {
	s.Lock()
	ctx := stream.Context()
	p, ok := peer.FromContext(ctx)
	if !ok {
		s.Unlock()
		return nil
	}
	info, err := stream.Recv()
	if err != nil {
		s.Unlock()
		return err
	}
	switch v := p.Addr.(type) {
	case *net.TCPAddr:
		info.IP = v.IP.String()
	default:
		s.Unlock()
		return errors.New(v.Network() + " is not supported")
	}
	addr := fmt.Sprintf("%s:%d", info.IP, info.Port)
	streamInfo, ok := s.streams[addr]
	if ok {
		info = streamInfo.info
		streamInfo.isDown = false
	} else {
		if _, ok := s.ids[info.Type]; !ok {
			s.ids[info.Type] = 1
		}
		info.ID = s.ids[info.Type]
		info.TimeBase = s.timeBase
		s.ids[info.Type]++
		streamInfo = &StreamInfo{
			info:       info,
			serverDown: make(chan *rpc.ServerInfo, 10),
			serverUp:   make(chan *rpc.ServerInfo, 10),
		}
		s.streams[addr] = streamInfo
	}
	log.Println("register server", info)

	stream.Send(&rpc.RegisterReturn{
		ID:         info.ID,
		NewServers: s.getAllServers(),
	})

	for _, stream1 := range s.streams {
		if stream1 != streamInfo {
			stream1.serverUp <- info
		}
	}
	if _, ok := s.servers[info.Type]; !ok {
		s.servers[info.Type] = &rpc.ServerInfoIDMap{
			Servers: make(map[uint32]*rpc.ServerInfo),
		}
	}
	s.servers[info.Type].Servers[info.ID] = info
	s.Unlock()
	for {
		select {
		case up := <-streamInfo.serverUp:
			s.onServerUp(stream, up)
			break
		case down := <-streamInfo.serverDown:
			s.onServerDown(stream, down)
			break
		case <-ctx.Done():
			s.Lock()
			streamInfo.isDown = true
			for _, stream1 := range s.streams {
				if stream1 != streamInfo {
					stream1.serverDown <- info
				}
			}
			s.Unlock()
			err := ctx.Err()
			log.Println(ctx.Err())
			return err
		}
	}
	return nil
}

func (s *server) Sync(ctx context.Context, info *rpc.ServerInfo) (*rpc.RegisterReturn, error) {
	return nil, nil
}

func main() {
	var port = ""
	var timeBase int64
	flag.StringVar(&port, "port", "9999", "-port 9999")
	flag.Int64Var(&timeBase, "timebase", 0, "-timebase 1534061219")
	flag.Parse()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return
	}

	server := grpc.NewServer()
	rpc.RegisterClusterServer(server, NewServer(timeBase))
	err = server.Serve(lis)
	log.Println(err)
}
