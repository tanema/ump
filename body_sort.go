package ump

import (
	"sort"
)

type (
	bodiesBy   func(b1, b2 *Body) bool
	bodySorter struct {
		bodies []*Body
		by     bodiesBy
	}
)

func (by bodiesBy) Sort(bodies []*Body) {
	ps := &bodySorter{
		bodies: bodies,
		by:     by,
	}
	sort.Sort(ps)
}

func (s *bodySorter) Len() int {
	return len(s.bodies)
}

func (s *bodySorter) Swap(i, j int) {
	s.bodies[i], s.bodies[j] = s.bodies[j], s.bodies[i]
}

func (s *bodySorter) Less(i, j int) bool {
	return s.by(s.bodies[i], s.bodies[j])
}
