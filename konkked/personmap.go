package main

import (
	"sort"
)

type By func(p1, p2 *Person) bool

type Person struct {
	name  string
	stars int
}

type peopleSorter struct {
	people []*Person
	by     By
}

type PersonMap map[string]*Person

func (by By) Sort(people []*Person) {
	ps := &peopleSorter{people: people, by: by}

	sort.Sort(ps)
}

func (s *peopleSorter) Len() int {
	return len(s.people)
}

func (s *peopleSorter) Swap(i, j int) {
	s.people[i], s.people[j] = s.people[j], s.people[i]
}

func (s *peopleSorter) Less(i, j int) bool {
	return s.by(s.people[i], s.people[j])
}

func MakePersonMap() PersonMap {
	retv := make(map[string]*Person)
	var vl = PersonMap(retv)
	return vl
}
