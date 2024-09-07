package serial

import "fmt"

func PrintConsole(instruction string) (int, error) {
	fmt.Println(instruction)
	return len(instruction), nil
}
