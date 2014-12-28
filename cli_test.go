package mycalc

import (
	"bytes"
	"testing"
)

type runTest struct {
	name  string
	input string
	value string
}

var runTests = []runTest{
	{"single", "1", "1\n"},
	{"add", "1+1", "2\n"},
	{"sub", "1-1", "0\n"},
	{"mul", "2*3", "6\n"},
	{"div", "4/2", "2\n"},
	{"addmul", "4+2*5", "14\n"},
	{"subdiv", "6-4/2", "4\n"},
	{"twoline", "1+2\n3+4", "3\n7\n"},
	{"op", "+", "unexpected item\n"},
	{"op and success", "+\n1+2", "unexpected item\n3\n"},
	{"not expected symbol", "&", "unexpected item\n"},
}

func TestRun(t *testing.T) {
	for _, test := range runTests {
		b := bytes.Buffer{}
		RunString(test.input, &b)
		sb := b.String()
		if sb != test.value {
			t.Errorf("%s: got\n\t%v\nexpeected\n\t%v", test.name, sb, test.value)
		}
	}
}
