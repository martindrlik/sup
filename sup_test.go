package main

import (
	"strings"
	"time"
)

func Example() {
	now = func() func() time.Time {
		t := time.Date(1, 1, 1, 1, 0, 0, 0, time.UTC)
		return func() time.Time {
			t = t.Add(30 * time.Minute)
			return t
		}
	}()
	r := strings.NewReader(`add 1h readin book
resume -1
start watching netflix
resume 1
ps
fixname 0 reading book
ps ^r
resume 2
add x y
resume *
ps (
dump
`)
	proc(r)
	// Output:
	// 	0 1h30m0s readin book
	//  1 1h0m0s watching netflix
	//  0 1h30m0s reading book
	// add 1h30m0s reading book
	// add 1h30m0s watching netflix
}
