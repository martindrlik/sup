package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	Active int
	Start  time.Time
	Tasks  []Task
}

func (p *Project) proc(cmd string, args []string) error {
	switch cmd {
	case "new":
		p.add(strings.Join(args, " "))
	case "start":
		i, err := parse(args[0])
		if err != nil {
			return err
		}
		if i < 0 {
			i = len(p.Tasks) + i
		}
		if l := len(p.Tasks); l == 0 || i < 0 || i > l {
			return errors.New("index out of range")
		}
		now := time.Now()
		if p.Active != NoneActive {
			t := p.Tasks[p.Active]
			t.Took += now.Sub(p.Start)
			p.Tasks[p.Active] = t
		}
		p.Active = i
		p.Start = now
	case "list":
		return p.list(os.Stdout)
	}
	return nil
}

func parse(s string) (int, error) {
	n, err := strconv.ParseInt(s, 36, 32)
	return int(n), err
}

func loadProject(isCreate bool, name string) (*Project, error) {
	if isCreate {
		return &Project{
			Active: NoneActive,
		}, nil
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	p := &Project{}
	err = d.Decode(p)
	return p, err
}

func (p *Project) add(name string) {
	p.Tasks = append(p.Tasks, Task{Name: name})
}

func (p *Project) list(w io.Writer) error {
	imax, nmax := 0, 0
	for i, t := range p.Tasks {
		if il := len(strconv.Itoa(i)); il > imax {
			imax = il
		}
		if nl := len(t.Name); nl > nmax {
			nmax = nl
		}
	}
	for i, t := range p.Tasks {
		si := strconv.Itoa(i)
		ni, nn := 1, 1
		ni += imax - len(si)
		nn += nmax - len(t.Name)
		fmt.Fprint(w, si)
		fmt.Fprint(w, strings.Repeat(" ", ni))
		fmt.Fprint(w, t.Name)
		fmt.Fprint(w, strings.Repeat(" ", nn))
		fmt.Fprintln(w, prettyDuration(t.Took))
	}
	return nil
}

func prettyDuration(d time.Duration) string {
	s := d.String()
	t := strings.TrimSuffix(s, "h0m0s")
	if s == t {
		return s
	}
	if d.Hours() > 1 {
		return s + " hours"
	}
	return s + " hour"
}

func (p *Project) save(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}
