package main

import (
	"hash/fnv"
	"strconv"
	"sync"
)

type ServiceRegistrationsList struct {
	sync.Mutex
	IsDirty bool
	Items   []ServiceRegistration
}

type ServiceRegistration struct {
	Service      string
	Instance     string
	PodNs        string
	PodName      string
	PodNode      string
	PodIP        string
	Port         int32
	Registration string
}

func (s *ServiceRegistration) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.Service))
	h.Write([]byte(s.Instance))
	h.Write([]byte(s.PodNs))
	h.Write([]byte(s.PodName))
	h.Write([]byte(s.PodNode))
	h.Write([]byte([]byte(strconv.Itoa(int(s.Port)))))
	h.Write([]byte(s.PodIP))
	h.Write([]byte(s.Registration))
	return h.Sum64()
}

// compare registration slices ignoring order
func registrationsEqual(x, y []ServiceRegistration) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[uint64]int, len(x))
	for _, _x := range x {
		diff[_x.Hash()]++
	}
	for _, _y := range y {
		h := _y.Hash()
		if _, ok := diff[h]; !ok {
			return false
		}
		diff[_y.Hash()] -= 1
		if diff[h] == 0 {
			delete(diff, h)
		}
	}
	return len(diff) == 0
}
