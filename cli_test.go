package mycalc

import (
	"bytes"
	"testing"
)

type lineTest struct {
	name  string
	input string
	value string
}

var lineTests = []lineTest{
	{"single", "1", "1"},
	{"add", "1+1", "2"},
	{"sub", "1-1", "0"},
	{"mul", "2*3", "6"},
	{"div", "4/2", "2"},
	{"addmul", "4+2*5", "14"},
	{"subdiv", "6-4/2", "4"},
}

func TestLine(t *testing.T) {
	for _, test := range lineTests {
		b := bytes.Buffer{}
		RunLine(test.input, &b)
		sb := b.String()
		if sb != test.value {
			t.Errorf("%s: got\n\t%v\nexpeected\n\t%v", test.name, sb, test.value)
		}
	}
}
