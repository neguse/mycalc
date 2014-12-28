package mycalc

import (
	"bytes"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var origin = pos{line: 1, col: 1}

var tEof1 = item{itemEof, origin, ""}
var tEof2 = item{itemEof, pos{line: 1, col: 2}, ""}
var tEof3 = item{itemEof, pos{line: 1, col: 3}, ""}
var tEof4 = item{itemEof, pos{line: 1, col: 4}, ""}

var lexTests = []lexTest{
	{"empty", "", []item{tEof1}},
	{"number1", "1", []item{item{itemDoubleLiteral, origin, "1"}, tEof2}},
	{"number11", "11", []item{item{itemDoubleLiteral, origin, "11"}, tEof3}},
	{"number0", "0", []item{item{itemDoubleLiteral, origin, "0"}, tEof2}},
	{"number1.0", "1.0", []item{item{itemDoubleLiteral, origin, "1.0"}, tEof4}},
	{"number1.", "1.", []item{item{itemError, origin, "digit not appear next to dot"}}},
	{"number01", "01", []item{item{itemDoubleLiteral, origin, "01"}, tEof3}},
	{".", ".", []item{item{itemError, origin, "unexpected character"}}},
	{"i", "i", []item{item{itemError, origin, "unexpected character"}}},
	{"two line", "1\n2", []item{
		item{itemDoubleLiteral, origin, "1"},
		item{itemEol, pos{line: 1, col: 2}, "\n"},
		item{itemDoubleLiteral, pos{line: 2, col: 1}, "2"},
		item{itemEof, pos{line: 2, col: 2}, ""}}},
	{"add", "1+2", []item{
		item{itemDoubleLiteral, origin, "1"},
		item{itemAdd, pos{line: 1, col: 2}, "+"},
		item{itemDoubleLiteral, pos{line: 1, col: 3}, "2"},
		tEof4}},
	{"sub", "1-2", []item{
		item{itemDoubleLiteral, origin, "1"},
		item{itemSub, pos{line: 1, col: 2}, "-"},
		item{itemDoubleLiteral, pos{line: 1, col: 3}, "2"},
		tEof4}},
	{"mul", "1*2", []item{
		item{itemDoubleLiteral, origin, "1"},
		item{itemMul, pos{line: 1, col: 2}, "*"},
		item{itemDoubleLiteral, pos{line: 1, col: 3}, "2"},
		tEof4}},
	{"div", "1/2", []item{
		item{itemDoubleLiteral, origin, "1"},
		item{itemDiv, pos{line: 1, col: 2}, "/"},
		item{itemDoubleLiteral, pos{line: 1, col: 3}, "2"},
		tEof4}},
	{"paren", "(1)", []item{
		item{itemLParen, origin, "("},
		item{itemDoubleLiteral, pos{line: 1, col: 2}, "1"},
		item{itemRParen, pos{line: 1, col: 3}, ")"},
		tEof4}},
}

func collect(t *lexTest) (items []item) {
	buf := bytes.NewBufferString(t.input)
	l := lex(buf)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEof || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}

type stringTest struct {
	name  string
	input item
	str   string
}

var stringTests = []stringTest{
	{"add", item{itemAdd, origin, "+"}, `add:"+"(1,1)`},
	{"sub", item{itemSub, origin, "-"}, `sub:"-"(1,1)`},
	{"mul", item{itemMul, origin, "*"}, `mul:"*"(1,1)`},
	{"div", item{itemDiv, origin, "/"}, `div:"/"(1,1)`},
	{"doubleLiteral", item{itemDoubleLiteral, origin, "1.0"}, `doubleLiteral:"1.0"(1,1)`},
	{"eol", item{itemEol, origin, ""}, `eol:""(1,1)`},
	{"eof", item{itemEof, origin, ""}, `eof:""(1,1)`},
	{"error", item{itemError, origin, "error"}, `error:"error"(1,1)`},
}

func TestString(t *testing.T) {
	for _, test := range stringTests {
		if test.input.String() != test.str {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, test.input, test.str)
		}
	}
}
