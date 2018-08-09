package lib

import (
	"net"
	"sync"

	"google.golang.org/grpc"

	context "golang.org/x/net/context"
)

type SessionData struct {
	ID         ID
	ClientIP   net.IP
	ClientPort uint16
	LoginTime  uint64
	Client     ClientInfo
	Extra      map[string]interface{}
}

type ClientInfo struct {
	ClientType    uint16
	ClientVersion string
	ProtoType     string
	ProtoVersion  string
	OSType        uint16
	OSVersion     string
	DeviceType    string
	DeviceNumber  string
}

type Session struct {
	SessionData
	sync.RWMutex
	ctx       context.Context
	msgStream grpc.ServerStream
	m         *SessionManager
}

type SessionManager struct {
	sync.RWMutex
	index    uint32
	sessions map[ID]*Session
	cache    sync.Pool
}

func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[ID]*Session),
		cache: sync.Pool{
			New: func() interface{} {
				return &Session{}
			},
		},
	}
	return m
}

func (m *SessionManager) OpenSession(connectorID uint16, addr net.Addr) *Session {
	m.Lock()
	defer m.Unlock()
	var port uint16
	var ip net.IP
	var id ID
	switch v := addr.(type) {
	case *net.TCPAddr:
		port = uint16(v.Port)
		ip = v.IP
	}

	for {
		m.index++
		if GetBits(uint64(m.index), IndexBits, 0) == 0 {
			m.index = 1
		}
		index := m.index
		id = MakeSessionID(connectorID, index)
		if _, ok := m.sessions[id]; !ok {
			break
		}
	}
	s := m.cache.Get().(*Session)
	s.SessionData = SessionData{
		ID:         id,
		ClientIP:   ip,
		ClientPort: port,
	}
	s.m = m
	m.sessions[id] = s
	return s
}

func (m *SessionManager) GetSessionCount() int {
	m.Lock()
	defer m.Unlock()
	return len(m.sessions)
}

func (m *SessionManager) CloseSession(id ID) {
	m.Lock()
	defer m.Unlock()
	s, ok := m.sessions[id]
	s.Lock()
	defer s.Unlock()
	if ok {
		delete(m.sessions, id)
		s.ID = 0
		s.LoginTime = 0
		s.Client = ClientInfo{}
		s.Extra = nil
		m.cache.Put(s)
	}
}

func (s *Session) UpdateClientInfo(c *ClientInfo) {
	s.Lock()
	defer s.Unlock()
	s.SessionData.Client = *c
}

func (s *Session) SetExtra(key string, val interface{}) {
	s.Lock()
	defer s.Unlock()
	if s.Extra == nil {
		s.Extra = make(map[string]interface{})
	}
	s.Extra[key] = val
}

func (s *Session) Close() {
	s.m.CloseSession(s.ID)
}
