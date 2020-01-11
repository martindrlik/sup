package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	names   []string
	takes   []time.Duration
	active  = -1
	started time.Time
	now     = time.Now
)

func main() {
	proc(os.Stdin)
}

var commands = map[string]func(string){
	"start":   start,
	"fixname": indexFor(fixname),
	"resume":  indexFor(resume),
	"ps":      exprFor(ps),
	"add":     add,
	"dump":    dump,
}

func proc(r io.Reader) {
	b := bufio.NewReader(r)
	for {
		t, err := b.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "sup: error reading input: %v\n", err)
			return
		}
		t = strings.TrimSpace(t)
		u, v := split(t)
		if f, ok := commands[u]; ok {
			f(v)
		} else {
			fmt.Fprintf(os.Stderr, "sup: unknown command: %s\n", u)
		}
		if err == io.EOF {
			break
		}
	}
}

func add(t string) {
	u, v := split(t)
	d, err := time.ParseDuration(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sup: error parsing duration: %v\n", err)
		return
	}
	names = append(names, v)
	takes = append(takes, d)
}

func start(name string) {
	i := len(names)
	names = append(names, name)
	takes = append(takes, 0)
	resume(i, "")
}

func fixname(i int, name string) {
	names[i] = name
}

func resume(i int, _ string) {
	t := now()
	if active >= 0 {
		takes[active] += t.Sub(started)
	}
	started = t
	active = i
}

func ps(re *regexp.Regexp) {
	for i, name := range names {
		if !re.MatchString(name) {
			continue
		}
		d := duration(i)
		d = d.Truncate(time.Second)
		fmt.Printf("%2v %5v %s\n", index(i), d, name)
	}
}

func index(i int) string {
	return strconv.FormatInt(int64(i), 36)
}

func duration(i int) time.Duration {
	d := takes[i]
	if i == active {
		d += now().Sub(started)
	}
	return d
}

func indexFor(fn func(int, string)) func(string) {
	return func(s string) {
		u, v := split(s)
		n, err := strconv.ParseInt(u, 36, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sup: parse index error: %v\n", err)
			return
		}
		i := int(n)
		l := len(takes)
		if i < 0 {
			i = l + i
		}
		if i < 0 || i >= l {
			fmt.Fprintf(os.Stderr, "sup: no task found for %s\n", u)
			return
		}
		fn(i, v)
	}
}

func exprFor(fn func(re *regexp.Regexp)) func(string) {
	return func(expr string) {
		re, err := regexp.Compile(expr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sup: error compiling expr: %v\n", err)
			return
		}
		fn(re)
	}
}

func split(t string) (u, v string) {
	s := strings.SplitN(t, " ", 2)
	u = s[0]
	if len(s) == 2 {
		v = s[1]
	}
	return
}

func dump(string) {
	for i, name := range names {
		d := duration(i)
		fmt.Println("add", d, name)
	}
}
