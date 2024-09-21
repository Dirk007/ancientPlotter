package serial

import (
	"errors"
	"fmt"
	"strings"

	hwserial "go.bug.st/serial"
)

var _ Writer = &SerialWriter{}

const (
	DefaultCT630Baud = 9600
	DefaultCT630Bits = 8
)

type SerialWriter struct {
	port hwserial.Port
}

func GuessPortName() (string, error) {
	fmt.Println("Trying to guess serial port...  This may take a few seconds.")
	ports, err := hwserial.GetPortsList()
	if err != nil {
		fmt.Printf("Failed to get serial ports: %v\n", err)
		return "", err
	}

	fmt.Printf("Found ports: %+v\n", ports)

	for _, port := range ports {
		lower := strings.ToLower(port)
		if strings.Contains(lower, "tty") && (strings.Contains(lower, "usb") || strings.Contains(lower, "serial")) {
			return port, nil
		}
	}

	fmt.Println("No port found that seems to be a serial port. Please check your device and try again. (Hint: try using 'ls /dev/tty*' or 'ls /dev/cu.*' on")
	return "", errors.New("unable to guess the right serial port")
}

func TryNew(device string) (*SerialWriter, error) {
	mode := &hwserial.Mode{
		BaudRate: DefaultCT630Baud,
		Parity:   hwserial.NoParity,
		DataBits: DefaultCT630Bits,
		StopBits: hwserial.OneStopBit,
		InitialStatusBits: &hwserial.ModemOutputBits{
			RTS: true,
			DTR: true,
		},
	}
	port, err := hwserial.Open(device, mode)
	if err != nil {
		return nil, err
	}
	return &SerialWriter{
		port: port,
	}, nil
}

func (w *SerialWriter) Write(data string) (int, error) {
	return w.port.Write([]byte(data))
}
