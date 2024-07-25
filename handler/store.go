package handler

import (
	"sync"
)

type AppStore struct {
	mutex       sync.RWMutex
	SuccessData map[int]Task
	FailedData  map[int]Task
}

func InitStore() *AppStore {
	return &AppStore{
		mutex:       sync.RWMutex{},
		SuccessData: make(map[int]Task, 0),
		FailedData:  make(map[int]Task, 0),
	}
}

func (s *AppStore) LoadSuccess(task Task) {
	s.mutex.Lock()
	if _, ok := s.SuccessData[task.Id]; !ok {
		s.SuccessData[task.Id] = task
	}
	s.mutex.Unlock()
}

func (s *AppStore) LoadFail(task Task) {
	s.mutex.Lock()
	if _, ok := s.FailedData[task.Id]; !ok {
		s.FailedData[task.Id] = task
	}
	s.mutex.Unlock()
}

func (s *AppStore) ClearMap() {
	s.mutex.Lock()
	clear(s.FailedData)
	clear(s.SuccessData)
	s.mutex.Unlock()
}
