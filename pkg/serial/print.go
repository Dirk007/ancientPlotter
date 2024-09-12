package serial

import (
	"fmt"
	"time"
)

func PrintConsole(instruction string) (int, error) {
	fmt.Println(instruction)
	// FIXME: Remove sleep
	time.Sleep(4 * time.Millisecond)
	return len(instruction), nil
}
