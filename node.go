package mycalc

import (
	"errors"
	"fmt"
)

type value struct {
	v   float64
	err error
}

type node interface {
	Type() nodeType
	Evaluate() value
	String() string
}

type nodeType int

const (
	nodeOp nodeType = iota
	nodeValue
	nodeError
)

type opType int

const (
	opAdd opType = iota
	opSub
	opMul
	opDiv
)

type opNode struct {
	opTyp    opType
	lhs, rhs node
}

func newOpNode(lhs node, rhs node, op opType) *opNode {
	return &opNode{
		opTyp: op,
		lhs:   lhs,
		rhs:   rhs,
	}
}

func (n *opNode) Type() nodeType {
	return nodeOp
}

func (n *opNode) Evaluate() value {
	lv := n.lhs.Evaluate()
	rv := n.rhs.Evaluate()
	if lv.err != nil {
		return value{0, lv.err}
	} else if rv.err != nil {
		return value{0, rv.err}
	}
	switch n.opTyp {
	case opAdd:
		return value{lv.v + rv.v, nil}
	case opSub:
		return value{lv.v - rv.v, nil}
	case opMul:
		return value{lv.v * rv.v, nil}
	case opDiv:
		return value{lv.v / rv.v, nil}
	default:
		return value{0, errors.New("unexpected opTyp")}
	}
}

func (n *opNode) String() string {
	var opStr string
	switch n.opTyp {
	case opAdd:
		opStr = "+"
	case opSub:
		opStr = "-"
	case opMul:
		opStr = "*"
	case opDiv:
		opStr = "/"
	default:
		opStr = fmt.Sprintf("%d", n.opTyp)
	}
	return fmt.Sprintf("(%s %s %s)", opStr, n.lhs.String(), n.rhs.String())
}

type valueNode struct {
	v float64
}

func newValueNode(v float64) *valueNode {
	return &valueNode{
		v: v,
	}
}

func (n *valueNode) Type() nodeType {
	return nodeValue
}

func (n *valueNode) Evaluate() value {
	return value{n.v, nil}
}

func (n *valueNode) String() string {
	return fmt.Sprintf("%f", n.v)
}

type errorNode struct {
	err string
}

func newErrorNode(err string) *errorNode {
	return &errorNode{
		err: err,
	}
}

func (n *errorNode) Type() nodeType {
	return nodeError
}

func (n *errorNode) Evaluate() value {
	return value{0, errors.New(n.err)}
}

func (n *errorNode) String() string {
	return n.err
}
