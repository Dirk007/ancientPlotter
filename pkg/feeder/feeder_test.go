package feeder_test

import (
	"context"
	"testing"
	"time"

	"github.com/Dirk007/ancientPlotter/pkg/feeder"
	"github.com/Dirk007/ancientPlotter/pkg/serial"
	"github.com/stretchr/testify/mock"
)

var _ serial.Writer = &MockWriter{}

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(s string) (n int, err error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func TestResumesLine(t *testing.T) {
	ctx := context.Background()

	writer := new(MockWriter)
	writer.On("Write", "1234567890;").Return(5, nil)
	writer.On("Write", "67890;").Return(6, nil)

	f := feeder.New(2, time.Millisecond, "test-job", writer)
	f.WriteInstruction(ctx, "1234567890;", func(stats feeder.Stats) error { return nil }, &feeder.Stats{})

	writer.AssertExpectations(t)
}
