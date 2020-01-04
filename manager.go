package main

import (
	"fmt"
	"strconv"
	"time"
)

type manager struct {
	names []string
	tooks []time.Duration

	m map[string]int
	t time.Time
	c int
}

func newManager() *manager {
	return &manager{
		m: make(map[string]int),
	}
}

func (m *manager) start(name string) {
	t := time.Now()
	if len(m.tooks) > 0 {
		m.tooks[len(m.tooks)-1] += t.Sub(m.t)
	}
	m.t = t
	if i, ok := m.m[name]; ok {
		m.c = i
		return
	}
	m.names = append(m.names, name)
	m.tooks = append(m.tooks, 0)
	i := len(m.names) - 1
	m.m[name] = i
	m.c = i
}

func (m *manager) ps(pattern string) {
	const format = "%4v %5v %s\n"
	fmt.Printf(format, "id", "took", "name")
	var total time.Duration
	for i, name := range m.names {
		d := m.tooks[i]
		if i == m.c {
			d += time.Now().Sub(m.t)
		}
		total += d
		d = d.Truncate(config.td)
		fmt.Printf(format, i, d, name)
	}
	total = total.Truncate(config.td)
	fmt.Printf("     %5v\n", total)
}

func (m *manager) fixname(s string) {
	s, rest := firstWord(s)
	i, err := strconv.Atoi(s)
	notFound := func() {
		sup.Printf("task not found: %s", s)
	}
	if err != nil {
		notFound()
		return
	}
	l := len(m.names)
	if i < 0 {
		i = l + i
	}
	if i < 0 || i >= l {
		notFound()
		return
	}
	m.names[i] = rest
}
