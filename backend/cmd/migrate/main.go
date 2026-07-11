package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: migrate [up|down]")
		os.Exit(2)
	}
	fmt.Printf("migration runner boundary established; configure a production migration adapter before %s\n", os.Args[1])
}
