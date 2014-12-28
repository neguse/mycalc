package main

import (
	"os"

	"github.com/neguse/mycalc"
)

func main() {
	mycalc.RunReader(os.Stdin, os.Stdout)
}
