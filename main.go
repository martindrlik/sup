package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var config = struct {
	td time.Duration
}{
	td: time.Second,
}

var cmds = map[string]func(string){
	"start":   m.start,
	"ps":      m.ps,
	"fixname": m.fixname,

	"truncate": truncate,

	"help": help,
}

var (
	sup = log.New(os.Stdout, "sup: ", 0)
	m   = newManager()
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			sup.Fatal(err)
		}
		line = strings.TrimSpace(line)
		s, arg := firstWord(line)
		if cmd, ok := cmds[s]; ok {
			cmd(arg)
			continue
		}
		sup.Printf("unknown command: %s\n", s)
	}
}

func firstWord(s string) (first, rest string) {
	z := strings.SplitN(s, " ", 2)
	if len(z) == 1 {
		return z[0], ""
	}
	return z[0], z[1]
}

func help(string) {
	fmt.Println("start name   starts new or resumes existing task")
	fmt.Println("ps [pattern] prints tasks filtered by pattern")
	fmt.Println("fixname i name sets i-th task name to name")
	fmt.Println("truncate d (for printing) rounds durations toward zero to a multiple of d")
}

func truncate(s string) {
	d, err := time.ParseDuration(s)
	if err != nil {
		sup.Println(err)
		return
	}
	config.td = d
}
