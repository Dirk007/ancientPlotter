package serial

import (
	"fmt"
	"time"
)

var _ Writer = &ConsoleWriter{}

type ConsoleWriter struct{}

func (*ConsoleWriter) Write(instruction string) (int, error) {
	fmt.Println(instruction)
	// FIXME: Remove sleep
	time.Sleep(4 * time.Millisecond)
	return len(instruction), nil
}
