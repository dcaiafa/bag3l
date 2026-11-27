//go:build !bag3ldebug

package vm

// debugStackCheck enables a per-instruction assertion that the operand stack
// stays within the bound computed by computeMaxStack. Enable it with
// `-tags bag3ldebug`.
const debugStackCheck = false
