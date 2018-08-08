package lib

import (
	"net"
	"sync"
)

type Session struct {
	ID         ID
	Connection ConnectionInfo
	Client     ClientInfo
	Extra      map[string]string
}

type ClientInfo struct {
	LoginTime     uint64
	ClientType    uint16
	ClientVersion string
	ProtoType     string
	ProtoVersion  string
	OSType        uint16
	OSVersion     string
	DeviceType    string
	DeviceNumber  string
}

type ConnectionInfo struct {
	ConnectorID uint16
	LocalPort   uint16
	ClientIP    string
	ClientPort  uint16
}

type session struct {
	Session
	sync.RWMutex
}

type SessionManager struct {
	sync.RWMutex
	sessions map[ID]*session
	cache    sync.Pool
}

func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[ID]*session),
		cache: sync.Pool{
			New: func() interface{} {
				return &session{}
			},
		},
	}
	return m
}

func (m *SessionManager) OpenSession(connectorID uint16, addr net.Addr) *Session {
	m.Lock()
	defer m.Unlock()
	var localPort uint16
	switch v := addr.(type) {
	case *net.TCPAddr:
		localPort = uint16(v.Port)
	}
	id := MakeSessionID(connectorID, localPort)
	s := &session{}
	s.Session = Session{
		ID: id,
		Connection: ConnectionInfo{
			ConnectorID: connectorID,
			LocalPort:   localPort,
		},
	}
	m.sessions[id] = s
	return &s.Session
}

func (m *SessionManager) CloseSession(id ID) {
	m.Lock()
	defer m.Unlock()
	s, ok := m.sessions[id]
	if ok {
		s.ID = 0
		m.cache.Put(s)
	}
}
