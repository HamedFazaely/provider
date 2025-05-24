package db

import (
	"sync"
	"time"
)

type OrderStatus struct {
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type StatsDB struct {
	m  map[uint64][]OrderStatus
	mu sync.RWMutex
}

func NewStatsDB() *StatsDB {
	return &StatsDB{
		m: make(map[uint64][]OrderStatus),
	}
}

func (s *StatsDB) UpdateStatus(id uint64, stat *OrderStatus) {
	s.mu.Lock()
	stat.CreatedAt = time.Now()
	if _, ok := s.m[id]; !ok {
		s.m[id] = make([]OrderStatus, 0)
	}
	s.m[id] = append(s.m[id], *stat)
	s.mu.Unlock()
}

func (s *StatsDB) GetOrderStatuses(id uint64) []OrderStatus {
	var ret []OrderStatus
	s.mu.RLock()
	ret = append(ret, s.m[id]...)
	s.mu.RUnlock()
	return ret
}
