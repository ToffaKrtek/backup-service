package socket

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockConnHandler struct {
	mock.Mock
}

func (m *MockConnHandler) Handle(conn net.Conn) {
	log.Print("вызов хендлера")
	m.Called(conn)
}

func TestSocketStart(t *testing.T) {
	mockHandler := new(MockConnHandler)

	mockHandler.On("Handle", mock.Anything).Times(3)
	go SocketStart(mockHandler.Handle, mockHandler.Handle, mockHandler.Handle)

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("unix", socketPath)

	if err != nil {
		t.Fatalf("Не удалось подключится к сокету: %v", err)
	}
	defer conn.Close()
	time.Sleep(100 * time.Millisecond)
	mockHandler.AssertExpectations(t)
}
