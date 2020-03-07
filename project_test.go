package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestLoadProject(t *testing.T) {
	p, err := loadProject(false, "./testdata/p")
	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
	if p.Active != 1 {
		t.Errorf("expected task 1 to be active got %v", p.Active)
	}
	if len(p.Tasks) != 2 {
		t.Errorf("expected two tasks got %v", len(p.Tasks))
	}
	if p.Tasks[0].Name != "playing games" {
		t.Errorf("expected task to be named \"playing games\" got %v", p.Tasks[0].Name)
	}
	if p.Tasks[0].Took != 2*time.Hour {
		t.Errorf("expected task to take two hours got %v", p.Tasks[0].Took)
	}
	if p.Tasks[1].Name != "reading books" {
		t.Errorf("expected task to be name \"reading books\" got %v", p.Tasks[1].Name)
	}
	if p.Tasks[1].Took != time.Hour {
		t.Errorf("expected task to take two hours got %v", p.Tasks[1].Took)
	}
}

func TestCreateProject(t *testing.T) {
	p, err := loadProject(true, "./testadata/p")
	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
	if p.Active != NoneActive {
		t.Errorf("expected no active task got %v", p.Active)
	}
}

func TestSave(t *testing.T) {
	p := &Project{
		Active: 2,
	}
	w := &bytes.Buffer{}
	p.save(w)
	if w.Len() == 0 {
		t.Error("expected to write some bytes")
	}
}

func TestProcNew(t *testing.T) {
	p := &Project{}
	p.proc("new", []string{"playing"})
	if p.Tasks[0].Name != "playing" {
		t.Errorf("expected adding task with name \"playing\" got %v", p.Tasks[0].Name)
	}
}

func TestProcStartParseError(t *testing.T) {
	p := &Project{}
	err := p.proc("start", []string{"+"})
	msg := `strconv.ParseInt: parsing "+": invalid syntax`
	if err.Error() != msg {
		t.Errorf("expected parsing error %v got %v", msg, err)
	}
}

func TestProcStart9(t *testing.T) {
	p := &Project{}
	err := p.proc("start", []string{"9"})
	msg := "index out of range"
	if err.Error() != msg {
		t.Errorf("expected parsing error %v got %v", msg, err)
	}
}

func TestProcStartNegative(t *testing.T) {
	p := &Project{}
	p.Tasks = []Task{
		Task{Name: "playing"},
	}
	err := p.proc("start", []string{"-1"})
	if err != nil {
		t.Fatal(err)
	}
	if p.Active != 0 {
		t.Errorf("-1 should active last task (0) got %v", p.Active)
	}
}

func TestProcList(t *testing.T) {
	p := &Project{}
	p.Tasks = []Task{
		Task{Name: "playing games", Took: time.Hour},
		Task{Name: "playing with go", Took: 2 * time.Hour},
	}
	err := p.proc("list", []string{"-1"})
	if err != nil {
		t.Fatal(err)
	}
	expect := &bytes.Buffer{}
	fmt.Fprintln(expect, "0 playing games   1 hour")
	fmt.Fprintln(expect, "1 playing with go 2 hours")
	actual := &bytes.Buffer{}
	err = p.list(actual)
	if err != nil {
		t.Fatal(err)
	}
	if expectString, actualString := expect.String(), actual.String(); expectString != actualString {
		t.Errorf("expected\n%vgot\n%v", expectString, actualString)
	}
}
