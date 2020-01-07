package main

import (
	"log"
	"regexp"
)

func filter(expr string) (res []Task) {
	if expr == "" {
		return all
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		log.Println(err)
		return
	}
	for _, task := range all {
		if !re.MatchString(task.Name) {
			continue
		}
		res = append(res, task)
	}
	return
}
