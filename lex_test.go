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

var tEof = item{itemEof, 0, ""}

var lexTests = []lexTest{
	{"empty", "", []item{tEof}},
	{"number1", "1", []item{item{itemDoubleLiteral, 0, "1"}, tEof}},
	{"number11", "11", []item{item{itemDoubleLiteral, 0, "11"}, tEof}},
	{"number0", "0", []item{item{itemDoubleLiteral, 0, "0"}, tEof}},
	{"number1.0", "1.0", []item{item{itemDoubleLiteral, 0, "1.0"}, tEof}},
	{"number1.", "1.", []item{item{itemError, 0, "digit not appear next to dot"}}},
	{"two line", "1\n2", []item{
		item{itemDoubleLiteral, 0, "1"},
		item{itemEol, 1, "\n"},
		item{itemDoubleLiteral, 2, "2"},
		tEof}},
	{"add", "1+2", []item{
		item{itemDoubleLiteral, 0, "1"},
		item{itemAdd, 1, "+"},
		item{itemDoubleLiteral, 2, "2"},
		tEof}},
	{"sub", "1-2", []item{
		item{itemDoubleLiteral, 0, "1"},
		item{itemSub, 1, "-"},
		item{itemDoubleLiteral, 2, "2"},
		tEof}},
	{"mul", "1*2", []item{
		item{itemDoubleLiteral, 0, "1"},
		item{itemMul, 1, "*"},
		item{itemDoubleLiteral, 2, "2"},
		tEof}},
	{"div", "1/2", []item{
		item{itemDoubleLiteral, 0, "1"},
		item{itemDiv, 1, "/"},
		item{itemDoubleLiteral, 2, "2"},
		tEof}},
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

func equal(i1, i2 []item, checkPos bool) bool {
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
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
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
	{"add", item{itemAdd, 0, "+"}, "add:\"+\""},
	{"sub", item{itemSub, 0, "-"}, "sub:\"-\""},
	{"mul", item{itemMul, 0, "*"}, "mul:\"*\""},
	{"div", item{itemDiv, 0, "/"}, "div:\"/\""},
	{"doubleLiteral", item{itemDoubleLiteral, 0, "1.0"}, "doubleLiteral:\"1.0\""},
	{"eol", item{itemEol, 0, ""}, "eol:\"\""},
	{"eof", item{itemEof, 0, ""}, "eof"},
	{"error", item{itemError, 0, "error"}, "error"},
}

func TestString(t *testing.T) {
	for _, test := range stringTests {
		if test.input.String() != test.str {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, test.input, test.str)
		}
	}
}
