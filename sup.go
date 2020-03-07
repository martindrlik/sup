package main

import (
	"flag"
	"fmt"
	"os"
)

const NoneActive = -1

var (
	isCreate = flag.Bool("create", false, "create new project with name")
	projFile = flag.String("file", "sup-project", "project's file")
)

func main() {
	flag.Parse()
	p, err := loadProject(*isCreate, *projFile)
	fatal("unable to load project:", err)
	cmd, args := flag.Arg(0), arguments()
	fatal("unable to process command:", p.proc(cmd, args))
	f, err := os.Create(*projFile)
	fatal("unable to create project file:", err)
	fatal("unable to save project:", p.save(f))
}

func arguments() []string {
	if flag.NArg() > 1 {
		return flag.Args()[1:]
	}
	return nil
}

func fatal(prefix string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "sup:", prefix, err)
		os.Exit(1)
	}
}
