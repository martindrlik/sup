package main

import "strings"

var commands = map[string]command{
	"start": {
		usage:   "start \"reading book\"",
		isValid: func(s []string) bool { return len(s) > 1 },
		execute: func(s []string) { startTask(strings.Join(s[1:], " ")) },
	},
	"ls": {
		execute: func(s []string) { listTasks() },
	},
}

type command struct {
	usage   string
	isValid func([]string) bool
	execute func([]string)
}
