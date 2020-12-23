package set
/**
 * @DateTime   : 2020/12/23
 * @Author     : xumamba
 * @Description:
 **/

import (
	"sync"
)


type iSet interface {
	Add(items... interface{})
	Remove(items... interface{})
	Exists(item interface{}) bool
	Len() int
	List() []interface{}
	Clear()
}

type set struct {
	locker *sync.RWMutex

	container map[interface{}]struct{}
	len int
}

func (s *set) Add(items ...interface{}) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, item := range items{
		if !s.unsafeExists(item){
			s.container[item] = struct{}{}
			s.len++
		}
	}
}

func (s *set) Remove(items ...interface{}) {
	s.locker.Lock()
	defer s.locker.Unlock()
	for _, item := range items{
		if s.unsafeExists(item){
			delete(s.container, item)
			s.len--
		}
	}
}

func (s *set) unsafeExists(item interface{})bool  {
	_, ok := s.container[item]
	return ok
}

func (s *set) Exists(item interface{}) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.unsafeExists(item)
}

func (s *set) Len() int {
	s.locker.RLock()
	defer s.locker.RUnlock()
	return s.len
}

func (s *set) List() []interface{} {
	s.locker.RLock()
	defer s.locker.RUnlock()
	result := make([]interface{}, s.len)
	i := 0
	for elem := range s.container{
		result[i] = elem
		i++
	}
	return result
}

func (s *set) Clear() {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.len = 0
	s.container = map[interface{}]struct{}{}
}

func NewSet() *set {
	return &set{
		locker:    &sync.RWMutex{},
		container: map[interface{}]struct{}{},
		len:       0,
	}
}


