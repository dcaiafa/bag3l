package vm

type Fn struct {
	pkg       *CompiledPackage
	name      int
	locations []Location
	instrs    []Instr
	minArgs   int
	// maxStack is the peak operand-stack depth above the frame's locals.
	maxStack int
}

func (f *Fn) Type() string   { return "Func" }
func (f *Fn) String() string { return "<func>" }
func (f *Fn) Traits() Traits { return TraitNone }

func (f *Fn) Call(m *VM, args []Value, nRet int) ([]Value, error) {
	// Function calls are handled directly by the VM.
	panic("not called")
}
