package normalizer

import (
	"strings"

	"github.com/Dirk007/ancientPlotter/pkg/opcodes"
	"github.com/Dirk007/ancientPlotter/pkg/splitter"

	"github.com/sirupsen/logrus"
)

func DefaultCT630EndSequene() []string {
	return []string{
		"PU0,0;",   // Pen up and go home
		"SP0;SP0;", // Select NO pen, really
		" U F U @", // Some magic end sequence
	}
}

func (n *Normalizer) normalizeLine(line string) ([]string, error) {
	foundOpcode, err := opcodes.TryFrom(line)
	if err != nil {
		logrus.Warn(line)
		return nil, err
	}

	if foundOpcode != opcodes.OpcodePenDown {
		return []string{line}, nil
	}

	splitter, err := splitter.NewInstructionSplitter(line, ",")
	if err != nil {
		return nil, err
	}
	if n.replaceOpcode != nil && foundOpcode == n.replaceOpcode.Original {
		splitter = splitter.WithOpcodeReplacement(n.replaceOpcode.Replacement)
	}
	rebuilt, err := splitter.Rebuild()
	if err != nil {
		return nil, err
	}
	return rebuilt, nil
}

type OpcodeReplacement struct {
	Original    opcodes.Opcode
	Replacement opcodes.Opcode
}

type Normalizer struct {
	replaceOpcode *OpcodeReplacement
	endSequence   []string
}

func Default() *Normalizer {
	return &Normalizer{
		replaceOpcode: nil,
		endSequence:   DefaultCT630EndSequene(),
	}
}

func (n *Normalizer) WithOpcodeReplacement(from opcodes.Opcode, to opcodes.Opcode) *Normalizer {
	n.replaceOpcode = &OpcodeReplacement{
		Original:    from,
		Replacement: to,
	}
	return n
}

func (n *Normalizer) WithEndSequence(endSequence ...string) *Normalizer {
	n.endSequence = endSequence
	return n
}

func New(opcodeReplacement OpcodeReplacement, endSequence ...string) *Normalizer {
	return &Normalizer{
		replaceOpcode: &opcodeReplacement,
		endSequence:   endSequence,
	}
}

func (n *Normalizer) Normalize(input string) ([]string, error) {
	lines := strings.Split(input, ";")
	result := make([]string, 0)
	for index, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			logrus.Debugf("skipping empty input instruction at position %d/%d", index, len(lines))
			continue
		}
		normalized, err := n.normalizeLine(line)
		if err != nil {
			return nil, err
		}
		result = append(result, normalized...)
	}
	result = append(result, n.endSequence...)
	return result, nil
}
