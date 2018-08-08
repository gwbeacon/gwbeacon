package service

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"github.com/gwbeacon/sdk/v1"
	context "golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

type clientManager struct {
	sync.RWMutex
	clients map[string]*ClientState
}

type ClientState struct {
	addr    string
	user    *v1.UserInfo
	client  *v1.ClientInfo
	session *v1.Session
}

var manager *clientManager

func init() {
	manager = &clientManager{
		clients: make(map[string]*ClientState),
	}
}

func addClient(cs *ClientState) {
	manager.Lock()
	defer manager.Unlock()
	manager.clients[cs.addr] = cs
}

func getClient(addr string) *ClientState {
	manager.RLock()
	defer manager.RUnlock()
	return manager.clients[addr]
}

func delClient(cs *ClientState) {
	delClientByAddr(cs.addr)
}

func delClientByAddr(addr string) {
	manager.Lock()
	defer manager.Unlock()
	delete(manager.clients, addr)
}

func MakeSession(cs *ClientState) *v1.Session {
	info := fmt.Sprintf("%s_%s_%s_%s", cs.addr, cs.user.User.Name, cs.user.User.Domain, cs.client.ClientType)
	result := md5.Sum([]byte(info))
	var sid = make([]byte, 16)
	for i := 0; i < 16; i++ {
		sid = append(sid, result[i])
	}
	session := &v1.Session{
		Sid:       string(sid),
		LoginTime: time.Now().Unix(),
	}
	cs.session = session
	return session
}

func NewClientState(ctx context.Context, user *v1.UserInfo) *ClientState {
	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr.String() == "" {
		return nil
	}
	cs := &ClientState{
		user: user,
		addr: p.Addr.String(),
	}
	addClient(cs)
	return cs
}

func GetClientState(ctx context.Context) *ClientState {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil
	}
	return getClient(p.Addr.String())
}
