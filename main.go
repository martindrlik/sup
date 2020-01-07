package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	now = time.Now
	sup = log.New(os.Stdout, "sup: ", 0)
)

var config = struct {
	truncate time.Duration
}{
	truncate: time.Second,
}

func main() {
	r := bufio.NewReader(os.Stdout)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			sup.Fatal(err)
		}
		s = strings.TrimSpace(s)
		cmd, args := split2(s)
		if proc, ok := commands[cmd]; ok {
			proc(args)
			continue
		}
		sup.Println("unknown:", s)
	}
}

var commands = map[string]func(string){
	"start":   start,
	"resume":  resume,
	"ps":      list,
	"fixname": fixname,
}

var (
	mapping = map[string]int{}

	all []Task

	current int
	started time.Time
)

func start(name string) {
	setStarted()
	if i, ok := mapping[name]; ok {
		current = i
		return
	}
	current = len(all)
	all = append(all, Task{Name: name})
}

func resume(s string) {
	i, ok := index(s)
	if !ok {
		sup.Printf("no task found: %s", s)
		return
	}
	setStarted()
	current = i
}

func index(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	if i < 0 {
		i = len(all) + i
	}
	if i < 0 || i >= len(all) {
		return 0, false
	}
	return i, true
}

func setStarted() {
	t := now()
	if len(all) > 0 {
		task := all[current]
		task.Took += t.Sub(started)
		all[current] = task
	}
	started = t
}

func list(expr string) {
	s := filter(expr)
	for i, task := range s {
		d := task.Took
		if i == current {
			d += now().Sub(started)
		}
		d = d.Truncate(config.truncate)
		fmt.Printf("%2d %7v %s\n", i, d, task.Name)
	}
}

func fixname(args string) {
	s, name := split2(args)
	i, ok := index(s)
	if !ok {
		sup.Printf("no task found: %s", args)
		return
	}
	task := all[i]
	task.Name = name
	all[i] = task
}

func split2(s string) (string, string) {
	r := strings.SplitN(s, " ", 2)
	if len(r) == 1 {
		return r[0], ""
	}
	return r[0], r[1]
}
