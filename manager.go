package main

import (
	"fmt"
	"regexp"
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
	r, err := func() (*regexp.Regexp, error) {
		if pattern == "" {
			return nil, nil
		}
		r, err := regexp.Compile(pattern)
		return r, err
	}()
	if err != nil {
		sup.Println(err)
		return
	}
	fmt.Println("  id      took name")
	var total time.Duration
	for i, name := range m.names {
		if r != nil && !r.MatchString(name) {
			continue
		}
		d := m.tooks[i]
		if i == m.c {
			d += time.Now().Sub(m.t)
		}
		total += d
		d = d.Truncate(config.td)
		fmt.Printf("%4d %9v %s\n", i, d, name)
	}
	total = total.Truncate(config.td)
	fmt.Printf("     %9v\n", total)
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
