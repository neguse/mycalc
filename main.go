package main

import (
	"os"

	"github.com/neguse/mycalc/calc"
)

func main() {
	mycalc.RunReader(os.Stdin, os.Stdout)
}
