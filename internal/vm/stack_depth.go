package vm

import "fmt"

// maxOperandDepthLimit bounds the depth computeMaxStack is willing to track.
// Exceeding it means the bytecode grows the operand stack inside a loop, which
// is a compiler bug.
const maxOperandDepthLimit = 1 << 20

// computeMaxStack returns the peak number of operand-stack slots fn can occupy
// above its locals (i.e. above bp+nLocals), including slots that are only used
// transiently while an instruction executes.
//
// OpInitCallFrame uses the result to check for stack overflow once per call,
// which is what allows every push opcode to store to stack[sp] unchecked.
//
// The one thing that cannot be accounted for statically is the argument
// expansion performed by OpCall with CallExpandFlag; that is checked at run
// time.
func computeMaxStack(fn *Fn) int {
	instrs := fn.instrs
	if len(instrs) == 0 {
		return 0
	}

	// depthAt[ip] is the largest depth known on entry to instrs[ip], or -1 if
	// the instruction has not been reached. Where paths join with different
	// depths the larger one wins, which keeps the result an upper bound.
	depthAt := make([]int, len(instrs))
	for i := range depthAt {
		depthAt[i] = -1
	}

	var work []int
	visit := func(ip, depth int) {
		if ip >= len(instrs) || depth <= depthAt[ip] {
			return
		}
		if depth > maxOperandDepthLimit {
			panic(fmt.Sprintf("assert: unbounded operand stack at ip %v", ip))
		}
		depthAt[ip] = depth
		work = append(work, ip)
	}

	maxDepth := 0
	visit(0, 0)

	for len(work) > 0 {
		ip := work[len(work)-1]
		work = work[:len(work)-1]

		instr := instrs[ip]
		depth := depthAt[ip]

		// pop and push describe the instruction's net effect. peak is the
		// highest slot count, relative to the depth on entry, touched while
		// the instruction runs; it only needs setting when it exceeds the net
		// effect.
		pop, push, peak := 0, 0, 0
		falls := true // Execution can continue at ip+1.

		switch instr.opc {
		case OpNop, OpNot, OpUnaryMinus, OpInitCallFrame, OpMakeIter, OpSwap,
			OpLiftArg, OpInitLocal, OpInitLiftedLocal, OpInitGlobal:

		case OpJump:
			falls = false
			visit(int(instr.op1), depth)

		case OpJumpIfTrue, OpJumpIfFalse:
			pop = 1
			visit(int(instr.op1), depth-1)

		case OpDup, OpNil, OpNewInt, OpNewBool, OpNewMap, OpNewArray,
			OpLoadGlobal, OpLoadLocal, OpLoadLocalDeref, OpLoadArg,
			OpLoadArgDeref, OpLoadCapture, OpLoadCaptureBox, OpLoadLiteral:
			push = 1

		case OpPop:
			pop = int(instr.op1)

		case OpCall:
			// Pops the arguments and the callable; pushes the results.
			// Results are first written above the arguments (iterators) or
			// over them (native functions) and only then moved down over the
			// callable, hence the peak.
			nret := int(instr.op2)
			pop = int(instr.op1&CallArgCountMask) + 1
			push = nret
			peak = nret

		case OpNewClosure, OpNewIter:
			pop = int(instr.op2)
			push = 1

		case OpEvalBinOp, OpLoadIndex, OpArrayAppendNoPop,
			OpArrayExpandElemNoPop, OpDefer,
			OpStoreLocal, OpStoreLocalDeref, OpStoreArg, OpStoreArgDeref,
			OpStoreGlobal, OpStoreCapture:
			pop = 1

		case OpMapPutNoPop, OpSlice:
			pop = 2

		case OpStoreIndex:
			pop = 3

		case OpRet, OpIterRet:
			falls = false

		case OpThrow:
			pop = 1
			falls = false

		case OpBeginTry:
			// The handler is entered with the stack reset to the locals plus
			// the error value.
			visit(int(instr.op1), 1)

		case OpEndTry:
			falls = false
			visit(int(instr.op1), depth)

		case OpIterYield:
			// The yielded values are consumed by the caller; the iterator
			// resumes with them removed.
			pop = int(instr.op1)

		default:
			panic(fmt.Sprintf("assert: unknown opcode %v", instr.opc))
		}

		if pop > depth {
			panic(fmt.Sprintf(
				"assert: operand stack underflow at ip %v (opcode %v)",
				ip, instr.opc))
		}

		after := depth - pop + push
		maxDepth = max(maxDepth, after, depth+peak)

		if falls {
			visit(ip+1, after)
		}
	}

	return maxDepth
}
