package lib

import (
	"sync"
)

const (
	ContextSessionKey = "gwbeacon_session"
)

type Session struct {
	ID     ID
	Addr   string
	User   *UserInfo
	Client *ClientInfo
	Extra  map[string]interface{}
}

type UserInfo struct {
	Name      string
	Domain    string
	Device    string
	NickName  string
	LoginTime uint64
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

func NewSession(id ID, addr string) *Session {
	return &Session{
		ID:   id,
		Addr: addr,
	}
}

type SessionStore interface {
	Store(s *Session)
	Remove(id ID)
	Count() int
	Get(id ID) *Session
	Find(user *UserInfo) []*Session
}

type sessionStore struct {
	sync.RWMutex
	data  map[ID]*Session
	index map[string]map[string]map[string]*Session
}

func NewSessionStore() SessionStore {
	ss := &sessionStore{
		data:  make(map[ID]*Session),
		index: make(map[string]map[string]map[string]*Session),
	}
	return ss
}

func (ss *sessionStore) Store(s *Session) {
	ss.Lock()
	defer ss.Unlock()
	ss.data[s.ID] = s
	if s.User != nil && s.User.LoginTime > 0 {
		user := s.User
		if _, ok := ss.index[user.Domain]; !ok {
			ss.index[user.Domain] = make(map[string]map[string]*Session)
		}
		if _, ok := ss.index[user.Domain][user.Name]; !ok {
			ss.index[user.Domain][user.Name] = make(map[string]*Session)
		}
		ss.index[user.Domain][user.Name][user.Device] = s
	}
}

func (ss *sessionStore) Remove(id ID) {
	ss.Lock()
	defer ss.Unlock()
	if s, ok := ss.data[id]; ok && s.User != nil {
		user := s.User
		if _, ok := ss.index[user.Domain]; !ok {
			return
		}
		if _, ok := ss.index[user.Domain][user.Name]; !ok {
			return
		}
		if _, ok := ss.index[user.Domain][user.Name][user.Device]; ok {
			delete(ss.index[user.Domain][user.Name], user.Device)
		}
		if len(ss.index[user.Domain][user.Name]) == 0 {
			delete(ss.index[user.Domain], user.Name)
		}
		if len(ss.index[user.Domain]) == 0 {
			delete(ss.index, user.Domain)
		}
	}
}

func (ss *sessionStore) Count() int {
	ss.Lock()
	defer ss.Unlock()
	return len(ss.data)
}

func (ss *sessionStore) Get(id ID) *Session {
	ss.Lock()
	defer ss.Unlock()
	return ss.data[id]
}

func (ss *sessionStore) Find(user *UserInfo) []*Session {
	ss.Lock()
	defer ss.Unlock()
	if user.Domain == "" || user.Name == "" {
		return nil
	}
	if users, ok := ss.index[user.Domain]; ok {
		if devices, ok := users[user.Name]; ok {
			if user.Device != "" {
				if s, ok := devices[user.Device]; ok {
					return []*Session{s}
				}
			} else {
				var ret = make([]*Session, 0)
				for _, s := range devices {
					ret = append(ret, s)
				}
				return ret
			}
		}
	}
	return nil
}
