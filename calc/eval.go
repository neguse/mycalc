package mycalc

import (
	"errors"
	"fmt"
)

type env struct {
	global map[string]value
}

func newEnv() env {
	return env{
		global: make(map[string]value),
	}
}

func (e env) eval(n node) value {
	switch n.Type() {
	case nodeBinaryOp:
		nn := n.(*binaryOpNode)
		left := e.eval(nn.lhs)
		if left.err != nil {
			return left
		}
		right := e.eval(nn.rhs)
		if right.err != nil {
			return right
		}
		switch nn.opTyp {
		case binaryOpAdd:
			return value{v: left.v + right.v, err: nil}
		case binaryOpSub:
			return value{v: left.v - right.v, err: nil}
		case binaryOpMul:
			return value{v: left.v * right.v, err: nil}
		case binaryOpDiv:
			return value{v: left.v / right.v, err: nil}
		default:
			return value{v: 0, err: errors.New("unknown binaryOpType")}
		}
	case nodeUnaryOp:
		nn := n.(*unaryOpNode)
		v := e.eval(nn.operand)
		if v.err != nil {
			return v
		}
		switch nn.opTyp {
		case unaryOpMinus:
			return value{v: -v.v, err: nil}
		default:
			return value{v: 0, err: errors.New("unknown unaryOpType")}
		}
	case nodeValue:
		nn := n.(*valueNode)
		return value{v: nn.v, err: nil}
	case nodeError:
		nn := n.(*errorNode)
		return value{v: 0, err: errors.New(nn.err)}
	case nodeSushi:
		return value{v: 980, err: nil}
	case nodeAssign:
		nn := n.(*assignNode)
		v := e.eval(nn.expr)
		if v.err != nil {
			return v
		}
		e.global[nn.variable] = v
		return value{v: v.v, err: nil}
	case nodeVariableRef:
		nn := n.(*variableRefNode)
		v, ok := e.global[nn.variable]
		if !ok {
			return value{v: 0, err: errors.New(fmt.Sprint("unassigned variable ", nn.variable))}
		}
		return value{v: v.v, err: nil}
	default:
		return value{v: 0, err: errors.New("unknown type")}
	}
}
