package loginreq

import (
	"strings"
	"sync"
	"time"

	"github.com/rikkix/simplesso/utils/crypto"
)

type MemDB struct {
	dur time.Duration
	lock sync.RWMutex
	reqs map[string]LoginReq
	round uint32
	maxRound uint32
}

func NewMemDB(maxRound uint32, dur time.Duration) *MemDB {
	return &MemDB{
		dur: dur,
		lock: sync.RWMutex{},
		reqs: make(map[string]LoginReq),
		round: 0,
		maxRound: maxRound,
	}
}

func (m *MemDB) unsafeRemoveExpired() {
	for k, v := range m.reqs {
		if v.Expiry.Before(time.Now()) {
			delete(m.reqs, k)
		}
	}
}

func NewReqID() string {
	return crypto.RandomString(16)
}

func (m *MemDB) NewReq(username string, dur int) string {
	req := LoginReq{
		ID: NewReqID(),
		Username: username,
		Dur: dur,
		Expiry: time.Now().Add(m.dur),
		Confirmed: false,
		Code: "",
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.round++
	if m.round > m.maxRound {
		m.round = 0
		m.unsafeRemoveExpired()
	}
	m.reqs[req.ID] = req
	return req.ID
}

func (m *MemDB) Confirm(id string) (bool, string) {
	code := crypto.RandomDigits(8)

	m.lock.Lock()
	defer m.lock.Unlock()
	req, ok := m.reqs[id]
	if !ok {
		return false, ""
	}
	if req.Confirmed {
		return false, ""
	}
	req.Confirmed = true
	req.Code = code
	m.reqs[id] = req
	return true, code
}

func (m *MemDB) RemoveReq(id string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.reqs[id]
	if !ok {
		return false
	}
	delete(m.reqs, id)
	return true
}

func (m *MemDB) Finish(id string, code string) (bool, LoginReq) {
	m.lock.Lock()
	defer m.lock.Unlock()
	req, ok := m.reqs[id]
	if !ok {
		return false, LoginReq{}
	}
	if !req.Confirmed {
		return false, LoginReq{}
	}
	if !strings.EqualFold(req.Code, code) {
		delete(m.reqs, id)
		return false, LoginReq{}
	}
	delete(m.reqs, id)
	return true, req
}