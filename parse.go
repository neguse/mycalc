package mycalc

import (
	"fmt"
	"strconv"
)

type parser struct {
	input     chan item
	output    chan float64
	token     [1]item
	peekCount int
}

func (p *parser) nextValue() float64 {
	v := <-p.output
	return v
}

func parse(input chan item) *parser {
	p := &parser{
		input:  input,
		output: make(chan float64),
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
		for {
			t := p.peek().typ
			if t == itemEol || t == itemEof {
				break
			}
			p.next()
		}
		if p.peek().typ == itemEol {
			p.next()
		}
	}
	close(p.output)
}

func (p *parser) expression() node {
	t1 := p.term()
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
	if n.typ != itemDoubleLiteral {
		panic("unexpected item")
	}
	v, err := strconv.ParseFloat(n.val, 64)
	if err != nil {
		fmt.Print(v, err)
		panic("unexpected value")
	}
	return newValueNode(v)
}
