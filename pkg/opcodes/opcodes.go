package opcodes

import (
	"fmt"
	"strings"
)

type Opcode string

var _ fmt.Stringer = Opcode("foo")

// IN is always present, Could not found the actual definition of this.
// No idea if PA and PR are supported anyways.
const (
	OpcodeIN               Opcode = "IN"
	OpcodePenUp            Opcode = "PU"
	OpcodePenDown          Opcode = "PD"
	OpcodeSelectPen        Opcode = "SP"
	OpcodePositionAbsolute Opcode = "PA"
	OpcodePositionRelative Opcode = "PR"
)

func (op Opcode) String() string {
	return string(op)
}

func TryFrom(opCode string) (Opcode, error) {
	if len(opCode) < 2 {
		return "", fmt.Errorf("invalid opcode: '%s'", opCode)
	}
	opCode = opCode[:2]
	switch strings.ToUpper(opCode) {
	case string(OpcodeIN):
		return OpcodeIN, nil
	case string(OpcodePenUp):
		return OpcodePenUp, nil
	case string(OpcodePenDown):
		return OpcodePenDown, nil
	case string(OpcodeSelectPen):
		return OpcodeSelectPen, nil
	default:
		return "", fmt.Errorf("unknown opcode: %s", opCode)
	}
}
