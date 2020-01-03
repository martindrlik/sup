package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fs := strings.Fields(line)
		k := fs[0]
		c, ok := commands[k]
		if !ok {
			fmt.Println("\"start reading book\" starts task named \"reading book\"")
			fmt.Println("\"ls\" prints tasks")
			fmt.Printf("unknown: %q\n", k)
			continue
		}
		if c.isValid != nil && !c.isValid(fs) {
			fmt.Println("usage:", c.usage)
			continue
		}
		c.execute(fs)
	}
}
