package mycalc

type node interface {
	Type() nodeType
	Evaluate() float64
}

type nodeType int

const (
	nodeOp nodeType = iota
	nodeValue
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

func (n *opNode) Evaluate() float64 {
	lv := n.lhs.Evaluate()
	rv := n.rhs.Evaluate()
	switch n.opTyp {
	case opAdd:
		return lv + rv
	case opSub:
		return lv - rv
	case opMul:
		return lv * rv
	case opDiv:
		return lv / rv
	default:
		panic("unknown opTyp")
	}
}

type valueNode struct {
	value float64
}

func newValueNode(v float64) *valueNode {
	return &valueNode{
		value: v,
	}
}

func (n *valueNode) Type() nodeType {
	return nodeValue
}

func (n *valueNode) Evaluate() float64 {
	return n.value
}
