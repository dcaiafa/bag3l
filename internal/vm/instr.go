package vm

type Instr struct {
	opc OpCode
	op2 uint16
	op1 uint32
}

type OpCode byte

const (
	OpNop OpCode = iota
	OpJump
	OpJumpIfTrue
	OpJumpIfFalse
	OpDup
	OpPop
	OpCall
	OpNil
	OpNewClosure
	OpNewInt
	OpNewBool
	OpNewObject
	OpNewArray
	OpLoadGlobal
	OpLoadLocal
	OpLoadLocalDeref
	OpLoadArg
	OpLoadArgDeref
	OpLoadCapture
	OpLoadCaptureBox
	OpLoadLiteral
	OpEvalBinOp
	OpNot
	OpUnaryMinus
	OpObjectPutNoPop
	OpObjectGet
	OpArrayAppendNoPop
	OpArrayExpandElemNoPop
	OpRet
	OpStoreLocal
	OpStoreLocalDeref
	OpStoreArg
	OpStoreArgDeref
	OpStoreGlobal
	OpStoreCapture
	OpStoreIndex
	OpInitCallFrame
	OpMakeIter
	OpBeginTry
	OpEndTry
	OpSwap
	OpThrow
	OpDefer
	OpSlice
	OpIterYield
	OpIterRet
	OpNewIter
	OpLiftArg
	OpInitLocal
	OpInitLiftedLocal
	OpInitGlobal
)

const (
	CallArgCountMask uint32 = 0x7FFFFFFF
	CallExpandFlag   uint32 = 0x80000000

	OptionalIndexFlag uint16 = 0x0001
)
