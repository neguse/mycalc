package mycalc

import (
	"fmt"
	"io"
	"strings"
)

func RunString(line string, w io.Writer) {
	r := strings.NewReader(line)
	RunReader(r, w)
}

func RunReader(r io.Reader, w io.Writer) {
	l := lex(r)
	p := parse(l.items)
	for v := range p.output {
		if v.err == nil {
			fmt.Fprintln(w, v.v)
		} else {
			fmt.Fprintln(w, v.err)
		}
	}
}
