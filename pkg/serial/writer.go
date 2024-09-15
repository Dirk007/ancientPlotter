package serial

type Writer interface {
	Write(data string) (int, error)
}
