package splitter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dirk007/ancientPlotter/pkg/opcodes"
)

type InstructionSplitter struct {
	source        string
	sep           string
	splitted      []string
	opcode        string
	replaceOpcode *opcodes.Opcode
}

func NewInstructionSplitter(instruction string, sep string) (*InstructionSplitter, error) {
	if len(instruction) < 2 {
		return nil, fmt.Errorf("invalid instruction: '%s'", instruction)
	}
	opcode := instruction[:2]
	rest := instruction[2:]
	splitted := strings.Split(rest, sep)
	if len(splitted)%2 != 0 {
		return nil, fmt.Errorf("uneven instructions, need pairs. Have %d in '%s'", len(splitted), instruction)
	}

	return &InstructionSplitter{
		source:   instruction,
		sep:      sep,
		opcode:   opcode,
		splitted: splitted,
	}, nil
}

func (is *InstructionSplitter) WithOpcodeReplacement(opcode opcodes.Opcode) *InstructionSplitter {
	is.replaceOpcode = &opcode
	return is
}

func (is *InstructionSplitter) Take(n int) (string, bool) {
	if n > len(is.splitted) {
		return "", false
	}
	result := is.splitted[:n]
	is.splitted = is.splitted[n:]
	opcode := is.opcode
	if is.replaceOpcode != nil {
		opcode = (*is.replaceOpcode).String()
	}
	return fmt.Sprintf("%s%s%s%s", opcode, result[0], is.sep, result[1]), true
}

func (is *InstructionSplitter) Len() int {
	return len(is.splitted)
}

func (is *InstructionSplitter) Rebuild() ([]string, error) {
	result := make([]string, 0, is.Len())
	for {
		if is.Len() == 0 {
			break
		}
		command, ok := is.Take(2)
		if !ok {
			return nil, errors.New("uneven instructions, need pairs")
		}
		result = append(result, command)
	}
	return result, nil
}
