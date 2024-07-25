package handler

import "sync"

type AppStore struct {
	mutex       sync.Mutex
	SuccessData map[int]Task
	FailedData  map[int]Task
}

func InitStore() *AppStore {
	return &AppStore{
		mutex:       sync.Mutex{},
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

func (s *AppStore) CleanMap() {
	s.mutex.Lock()
	for key := range s.SuccessData {
		delete(s.SuccessData, key)
	}
	for key := range s.FailedData {
		delete(s.FailedData, key)
	}
	s.mutex.Unlock()
}
