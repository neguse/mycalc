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
	nodeBinaryOp nodeType = iota
	nodeUnaryOp
	nodeValue
	nodeError
)

type binaryOpType int

const (
	binaryOpAdd binaryOpType = iota
	binaryOpSub
	binaryOpMul
	binaryOpDiv
)

type binaryOpNode struct {
	opTyp    binaryOpType
	lhs, rhs node
}

func newBinaryOpNode(lhs node, rhs node, op binaryOpType) *binaryOpNode {
	return &binaryOpNode{
		opTyp: op,
		lhs:   lhs,
		rhs:   rhs,
	}
}

func (n *binaryOpNode) Type() nodeType {
	return nodeBinaryOp
}

func (n *binaryOpNode) Evaluate() value {
	lv := n.lhs.Evaluate()
	rv := n.rhs.Evaluate()
	if lv.err != nil {
		return value{0, lv.err}
	} else if rv.err != nil {
		return value{0, rv.err}
	}
	switch n.opTyp {
	case binaryOpAdd:
		return value{lv.v + rv.v, nil}
	case binaryOpSub:
		return value{lv.v - rv.v, nil}
	case binaryOpMul:
		return value{lv.v * rv.v, nil}
	case binaryOpDiv:
		return value{lv.v / rv.v, nil}
	default:
		return value{0, errors.New("unexpected opTyp")}
	}
}

func (n *binaryOpNode) String() string {
	var opStr string
	switch n.opTyp {
	case binaryOpAdd:
		opStr = "+"
	case binaryOpSub:
		opStr = "-"
	case binaryOpMul:
		opStr = "*"
	case binaryOpDiv:
		opStr = "/"
	default:
		opStr = fmt.Sprintf("%d", n.opTyp)
	}
	return fmt.Sprintf("(%s %s %s)", opStr, n.lhs.String(), n.rhs.String())
}

type unaryOpType int

const (
	unaryOpMinus unaryOpType = iota
)

type unaryOpNode struct {
	opTyp   unaryOpType
	operand node
}

func newUnaryOpNode(operand node, op unaryOpType) *unaryOpNode {
	return &unaryOpNode{
		opTyp:   op,
		operand: operand,
	}
}

func (n *unaryOpNode) Type() nodeType {
	return nodeUnaryOp
}

func (n *unaryOpNode) Evaluate() value {
	ov := n.operand.Evaluate()
	if ov.err != nil {
		return value{0, ov.err}
	}
	switch n.opTyp {
	case unaryOpMinus:
		return value{-ov.v, nil}
	default:
		return value{0, errors.New("unexpected opTyp")}
	}
}

func (n *unaryOpNode) String() string {
	var opStr string
	switch n.opTyp {
	case unaryOpMinus:
		opStr = "-"
	default:
		opStr = fmt.Sprintf("%d", n.opTyp)
	}
	return fmt.Sprintf("(%s %s)", opStr, n.operand.String())
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
