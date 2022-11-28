package store

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
	leaderWaitDelay     = 100 * time.Millisecond
	appliedWaitDelay    = 100 * time.Millisecond
)

type ConsistencyLevel int

const (
	Default    ConsistencyLevel = iota
	Stale                       // 允许读到旧数据
	Consistent                  // 强一致性读
)

// Store simple k-v store
type Store struct {
	mu sync.Mutex

	m map[string]string

	raft *raft.Raft // all changes are made by raft consense

	// Todo: logger

	RaftDir  string
	RaftBind string
}

func New() *Store {
	return &Store{
		m: make(map[string]string),
		// Todo:logger
	}
}

func (s *Store) LeaderID() (string, error) {
	addr := s.LeaderAddr()
	cfgFuture := s.raft.GetConfiguration()
	if err := cfgFuture.Error(); err != nil {
		return "", err
	}

	for _, srv := range cfgFuture.Configuration().Servers {
		if srv.Address == raft.ServerAddress(addr) {
			return string(srv.ID), nil
		}
	}

	return "", nil
}

func (s *Store) LeaderAddr() string {
	return string(s.raft.Leader())
}

// Join
func (s *Store) Join(nodeID string, httpAddr string, addr string) error {
	cfgFuture := s.raft.GetConfiguration()
	if err := cfgFuture.Error(); err != nil {
		return err
	}

	// 检查node是否已经存在
	for _, srv := range cfgFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) ||
			srv.Address == raft.ServerAddress(addr) {
			if srv.ID == raft.ServerID(nodeID) &&
				srv.Address == raft.ServerAddress(addr) {
				return nil
			}

			if err := s.raft.RemoveServer(srv.ID, 0, 0).Error(); err != nil {
				return fmt.Errorf("error removing existing node %v: %v", srv.ID, err)
			}

		}
	}

	if err := s.raft.AddNonvoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0).Error(); err != nil {
		return err
	}

	// 设置node基础信息
	if err := s.SetMeta(nodeID, addr); err != nil {
		return err
	}

	return nil
}

func (s *Store) SetMeta(key, value string) error {
	return s.Set(key, value)
}

func (s *Store) GetMeta(key string) (string, error) {
	return s.Get(key, Stale)
}

func (s *Store) LeaderAPIAddr() string {
	id, err := s.LeaderID()
	if err != nil {
		return ""
	}

	addr, err := s.GetMeta(id)
	if err != nil {
		return ""
	}

	return addr
}

// kv store item
type item struct {
	Op string `json:"op,omitempty"`

	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (s *Store) Set(key, value string) error {
	// only leader can write
	if s.raft.State() != raft.Leader {
		return raft.ErrNotLeader
	}

	i := item{Op: "set", Key: key, Value: value}

	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	// log copy
	return s.raft.Apply(b, raftTimeout).Error()
}

func (s *Store) Get(key string, lvl ConsistencyLevel) (string, error) {
	if lvl != Stale {
		if s.raft.State() != raft.Leader {
			return "", raft.ErrNotLeader
		}
	}

	if lvl == Consistent {
		if err := s.consistentRead(); err != nil {
			return "", err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key], nil
}

// consistentRead 通过验证leader是否变更
func (s *Store) consistentRead() error {
	if err := s.raft.VerifyLeader().Error(); err != nil {
		return err // fail fast if leader verification fails
	}

	return nil
}
