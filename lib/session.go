package lib

import (
	"sync"

	"errors"

	"github.com/gwbeacon/gwbeacon/lib/rpc"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	ContextSessionKey = "gwbeacon_session"
	ContextServerKey  = "gwbeacon_Server"
)

func NewSession(id uint64, addr string) *rpc.Session {
	return &rpc.Session{
		ID:   id,
		Addr: addr,
	}
}

type SessionStore interface {
	Save(s *rpc.Session) error
	Update(s *rpc.Session) error
	Remove(s *rpc.Session) ([]*rpc.Session, error)
	Replace(s *rpc.Session) ([]*rpc.Session, error)
	Stat() (*rpc.SessionStat, error)
	Get(s *rpc.Session) ([]*rpc.Session, error)
}

type sessionStoreClient struct {
	client rpc.SessionStoreClient
	cache  SessionStore
}

func NewSessionStoreClient(serverAddr string, withCache bool) (SessionStore, error) {
	conn, err := grpc.DialContext(context.Background(), serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := &sessionStoreClient{
		client: rpc.NewSessionStoreClient(conn),
	}
	if withCache {
		client.cache = NewSessionStore()
	}
	return client, nil
}

func (ss *sessionStoreClient) Save(s *rpc.Session) error {
	_, err := ss.client.Save(context.Background(), s)
	if err != nil {
		return err
	}
	if ss.cache != nil {
		ss.cache.Save(s)
	}
	return nil
}

func (ss *sessionStoreClient) Update(s *rpc.Session) error {
	_, err := ss.client.Update(context.Background(), s)
	if err != nil {
		return err
	}
	if ss.cache != nil {
		ss.cache.Update(s)
	}
	return nil
}

func (ss *sessionStoreClient) Remove(s *rpc.Session) ([]*rpc.Session, error) {
	res, err := ss.client.Remove(context.Background(), s)
	if err != nil {
		return nil, err
	}
	if ss.cache != nil {
		ss.cache.Remove(s)
	}
	return res.Data, nil
}

func (ss *sessionStoreClient) Replace(s *rpc.Session) ([]*rpc.Session, error) {
	res, err := ss.client.Replace(context.Background(), s)
	if err != nil {
		return nil, err
	}
	if ss.cache != nil {
		ss.cache.Replace(s)
	}
	return res.Data, nil
}

func (ss *sessionStoreClient) Stat() (*rpc.SessionStat, error) {
	return ss.client.Stat(context.Background(), &rpc.SessionStatRequest{})
}

func (ss *sessionStoreClient) Get(s *rpc.Session) ([]*rpc.Session, error) {
	if s.ID != 0 {
		if ss.cache != nil {
			res, _ := ss.cache.Get(s)
			if len(res) > 0 {
				return res, nil
			}
		}
	}
	res, err := ss.client.Get(context.Background(), s)
	return res.Data, err
}

type sessionStoreServer struct {
	sync.RWMutex
	cache map[uint64]SessionStore
}

func NewSessionStoreServer() SessionStore {
	return &sessionStoreServer{
		cache: make(map[uint64]SessionStore),
	}
}

func (ss *sessionStoreServer) getCache(sid uint64) SessionStore {
	ss.RLock()
	defer ss.RUnlock()
	connID := uint64(ID(sid).GetServerID())
	return ss.cache[connID]
}

func (ss *sessionStoreServer) getCacheByCreate(sid uint64) SessionStore {
	ss.Lock()
	defer ss.Unlock()
	connID := uint64(ID(sid).GetServerID())
	if cache, ok := ss.cache[connID]; ok {
		return cache
	}
	cache := NewSessionStore()
	ss.cache[connID] = cache
	return cache
}

func (ss *sessionStoreServer) Save(s *rpc.Session) error {
	cache := ss.getCacheByCreate(s.ID)
	return cache.Save(s)
}

func (ss *sessionStoreServer) Update(s *rpc.Session) error {
	cache := ss.getCacheByCreate(s.ID)
	return cache.Update(s)
}

func (ss *sessionStoreServer) Remove(s *rpc.Session) ([]*rpc.Session, error) {
	var ret = make([]*rpc.Session, 0)
	ss.Lock()
	var caches = make([]SessionStore, 0)
	for _, cache := range ss.cache {
		caches = append(caches, cache)
	}
	ss.Unlock()
	for _, cache := range caches {
		sessions, _ := cache.Remove(s)
		if len(sessions) > 0 {
			ret = append(ret, sessions...)
		}
	}
	if len(ret) > 0 {
		return ret, nil
	}
	return nil, errors.New("not found")
}

func (ss *sessionStoreServer) Replace(s *rpc.Session) ([]*rpc.Session, error) {
	s1 := &rpc.Session{
		User: s.User,
	}
	ret, _ := ss.Remove(s1)
	err := ss.Save(s)
	return ret, err
}

func (ss *sessionStoreServer) Stat() (*rpc.SessionStat, error) {
	ss.Lock()
	defer ss.Unlock()
	var ret = &rpc.SessionStat{
		Count:       0,
		DomainUsers: make(map[string]int32),
		ConnNumbers: make(map[uint64]int32),
	}
	for connID, cache := range ss.cache {
		stat, _ := cache.Stat()
		ret.Count += stat.Count
		if _, ok := ret.ConnNumbers[connID]; !ok {
			ret.ConnNumbers[connID] = 0
		}
		ret.ConnNumbers[connID] += stat.Count
		for domain, count := range stat.DomainUsers {
			if _, ok := ret.DomainUsers[domain]; !ok {
				ret.DomainUsers[domain] = 0
			}
			ret.DomainUsers[domain] += count
		}
	}
	return ret, nil
}

func (ss *sessionStoreServer) Get(s *rpc.Session) ([]*rpc.Session, error) {
	if s.ID != 0 {
		cache := ss.getCache(s.ID)
		return cache.Get(s)
	}
	return nil, errors.New("not found")
}

type sessionStore struct {
	sync.RWMutex
	data      map[uint64]*rpc.Session
	index     map[string]map[string]map[string]*rpc.Session
	indexAddr map[string]*rpc.Session
}

func NewSessionStore() SessionStore {
	ss := &sessionStore{
		data:      make(map[uint64]*rpc.Session),
		index:     make(map[string]map[string]map[string]*rpc.Session),
		indexAddr: make(map[string]*rpc.Session),
	}
	return ss
}

func (ss *sessionStore) Save(s *rpc.Session) error {
	ss.Lock()
	defer ss.Unlock()
	return ss.save(s)
}
func (ss *sessionStore) save(s *rpc.Session) error {
	ss.data[s.ID] = s
	ss.indexAddr[s.Addr] = s
	if s.User != nil && s.User.LoginTime > 0 {
		user := s.User
		if _, ok := ss.index[user.Domain]; !ok {
			ss.index[user.Domain] = make(map[string]map[string]*rpc.Session)
		}
		if _, ok := ss.index[user.Domain][user.Name]; !ok {
			ss.index[user.Domain][user.Name] = make(map[string]*rpc.Session)
		}
		ss.index[user.Domain][user.Name][user.Device] = s
	}
	return nil
}

func (ss *sessionStore) Update(s *rpc.Session) error {
	ss.Lock()
	defer ss.Unlock()
	if session, ok := ss.data[s.ID]; ok {
		if s.User != nil {
			session.User = s.User
		}
		if s.Client != nil {
			session.Client = s.Client
		}
		return nil
	}
	return errors.New("not found")
}

func (ss *sessionStore) Remove(s *rpc.Session) ([]*rpc.Session, error) {
	ss.Lock()
	defer ss.Unlock()
	if s.ID != 0 {
		return ss.removeByID(s.ID)
	} else if s.User != nil && s.User.LoginTime != 0 {
		return ss.removeByUser(s.User)
	}
	return nil, errors.New("not found")
}

func (ss *sessionStore) Replace(s *rpc.Session) ([]*rpc.Session, error) {
	ss.Lock()
	defer ss.Unlock()
	if s.User == nil || s.User.LoginTime == 0 || s.User.Device == "" {
		return nil, errors.New("wrong session")
	}
	ret, _ := ss.removeByUser(s.User)
	err := ss.save(s)
	return ret, err
}

func (ss *sessionStore) removeByUser(user *rpc.UserInfo) ([]*rpc.Session, error) {
	if users, ok := ss.index[user.Domain]; !ok {
		return nil, errors.New("not found")
	} else if devices, ok := users[user.Name]; ok {
		var ret = make([]*rpc.Session, 0)
		if user.Device == "" {
			for _, s := range devices {
				delete(ss.data, s.ID)
				delete(ss.indexAddr, s.Addr)
				ret = append(ret, s)
			}
			delete(ss.index[user.Domain], user.Name)
		} else if s, ok := devices[user.Device]; ok {
			delete(ss.data, s.ID)
			delete(ss.indexAddr, s.Addr)
			delete(devices, user.Device)
			ret = append(ret, s)
		}
		if len(devices) == 0 {
			delete(ss.index[user.Domain], user.Name)
		}
		if len(ss.index[user.Domain]) == 0 {
			delete(ss.index, user.Domain)
		}
		if len(ret) > 0 {
			return ret, nil
		}
	}

	return nil, errors.New("not found")
}

func (ss *sessionStore) removeByID(id uint64) ([]*rpc.Session, error) {
	s, ok := ss.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	delete(ss.data, id)
	if _, ok := ss.indexAddr[s.Addr]; ok {
		delete(ss.indexAddr, s.Addr)
	}
	if s.User != nil && s.User.LoginTime != 0 {
		ss.removeByUser(s.User)
	}
	return []*rpc.Session{s}, nil
}

func (ss *sessionStore) Stat() (*rpc.SessionStat, error) {
	ss.RLock()
	defer ss.RUnlock()
	ret := &rpc.SessionStat{
		Count:       int32(len(ss.data)),
		DomainUsers: make(map[string]int32),
	}
	for domain, users := range ss.index {
		ret.DomainUsers[domain] = int32(len(users))
	}
	return ret, nil
}

func (ss *sessionStore) Get(s *rpc.Session) ([]*rpc.Session, error) {
	ss.RLock()
	defer ss.RUnlock()
	if s.ID != 0 {
		if s1, ok := ss.data[s.ID]; ok {
			return []*rpc.Session{s1}, nil
		}
		return nil, errors.New("not found")
	} else if s.Addr != "" {
		if s1, ok := ss.indexAddr[s.Addr]; ok {
			return []*rpc.Session{s1}, nil
		}
		return nil, errors.New("not found")
	} else if s.User != nil && s.User.LoginTime != 0 {
		if users, ok := ss.index[s.User.Domain]; ok {
			if devices, ok := users[s.User.Name]; ok {
				if s.User.Device != "" {
					if s, ok := devices[s.User.Device]; ok {
						return []*rpc.Session{s}, nil
					}
				} else {
					var ret = make([]*rpc.Session, 0)
					for _, s := range devices {
						ret = append(ret, s)
					}
					return ret, nil
				}
			}
		}
	}
	return nil, errors.New("not found")
}
