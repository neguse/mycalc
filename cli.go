package mycalc

import (
	"fmt"
	"io"
	"strings"
)

func RunLine(line string, w io.Writer) {
	r := strings.NewReader(line)
	l := lex(r)
	p := parse(l.items)
	for v := range p.output {
		fmt.Fprint(w, v)
	}
}
