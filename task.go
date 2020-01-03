package main

import (
	"fmt"
	"time"
)

var (
	history = make([]string, 0, 20)
	tasks   = make(map[string]*Task)
)

type Task struct {
	d time.Duration
	t time.Time

	isStopped bool
}

func (u *Task) Start() {
	u.isStopped = false
	u.t = time.Now()
}

func (u *Task) Stop() {
	u.isStopped = true
	t := time.Now()
	u.d += t.Sub(u.t)
	u.t = t
}

func (u *Task) Duration() time.Duration {
	if u.isStopped {
		return u.d
	}
	return u.d + time.Now().Sub(u.t)
}

func startTask(k string) {
	if len(history) > 0 {
		last := history[len(history)-1]
		tasks[last].Stop()
	}
	u, ok := tasks[k]
	if !ok {
		u = &Task{}
		history = append(history, k)
		tasks[k] = u
	}
	u.Start()
}

func listTasks() {
	if len(history) < 1 {
		return
	}
	for _, k := range history {
		d := tasks[k].Duration().Truncate(time.Second)
		fmt.Println(k, d)
	}
}
