package mycalc

import (
	"errors"
	"strconv"
)

type parser struct {
	input     chan item
	output    chan value
	token     [1]item
	peekCount int
}

func (p *parser) nextValue() value {
	v, ok := <-p.output
	if ok {
		return v
	} else {
		return value{0, errors.New("receive failed")}
	}
}

func parse(input chan item) *parser {
	p := &parser{
		input:  input,
		output: make(chan value),
	}
	go p.parse()
	return p
}

func (p *parser) peek() item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = <-p.input
	return p.token[0]
}

func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = <-p.input
	}
	return p.token[p.peekCount]
}

func (p *parser) parse() {
	for p.peek().typ != itemEof {
		exp := p.expression()
		p.output <- exp.Evaluate()

		// Skip to next eol for recovering error.
		for p.peek().typ != itemEol && p.peek().typ != itemEof {
			// TODO: output as error
			p.next()
		}
		if p.peek().typ == itemEol {
			p.next()
		}
	}
	close(p.output)
}

func (p *parser) errorf(format string, args ...interface{}) {
}

func (p *parser) expression() node {
	t1 := p.term()
	if t1.Type() == nodeError {
		return t1
	}
	for {
		if p.peek().typ == itemAdd || p.peek().typ == itemSub {
			op := p.next()
			var t opType
			switch op.typ {
			case itemAdd:
				t = opAdd
			case itemSub:
				t = opSub
			default:
				panic("unknown itemType")
			}
			t2 := p.term()
			if t2.Type() == nodeError {
				return t2
			}
			// FIXME: Fix to left-to-right associative.
			t1 = newOpNode(t1, t2, t)
		} else {
			break
		}
	}
	return t1
}

func (p *parser) term() node {
	e1 := p.primaryExpression()
	if e1.Type() == nodeError {
		return e1
	}
	for {
		if p.peek().typ == itemMul || p.peek().typ == itemDiv {
			op := p.next()
			var t opType
			switch op.typ {
			case itemMul:
				t = opMul
			case itemDiv:
				t = opDiv
			default:
				panic("unknown itemType")
			}
			e2 := p.primaryExpression()
			if e2.Type() == nodeError {
				return e2
			}
			// FIXME: Fix to left-to-right associative.
			e1 = newOpNode(e1, e2, t)
		} else {
			break
		}
	}
	return e1
}

func (p *parser) primaryExpression() node {
	n := p.next()
	if n.typ == itemLParen {
		n2 := p.expression()
		if n2.Type() == nodeError {
			return n2
		}
		n3 := p.next()
		if n3.typ != itemRParen {
			return newErrorNode("expect RParen")
		}
		return n2
	} else if n.typ != itemDoubleLiteral {
		return newErrorNode("unexpected item")
	}
	v, err := strconv.ParseFloat(n.val, 64)
	if err != nil {
		return newErrorNode("unexpected value " + err.Error())
	}
	return newValueNode(v)
}
