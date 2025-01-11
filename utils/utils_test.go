package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set the working directory to the project root
	if err := os.Chdir("../"); err != nil {
		log.Fatalf("could not change working directory: %v", err)
	}

	// Run the tests
	os.Exit(m.Run())
}

// mockConn implements net.Conn interface for testing
type mockConn struct {
	readData  bytes.Buffer
	writeData bytes.Buffer
	closed    bool
	deadline  time.Time
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.closed {
		return 0, io.EOF
	}
	return m.readData.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	if m.closed {
		return 0, io.EOF
	}
	return m.writeData.Write(b)
}
func (m *mockConn) Close() error                       { m.closed = true; return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { m.deadline = t; return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// Helper function to create a new mock connection with input data
func newMockConn(input string) *mockConn {
	mock := &mockConn{}
	mock.readData.WriteString(input)
	return mock
}

// Helper function to read the written data from mock connection
func getWrittenData(mock *mockConn) string {
	return mock.writeData.String()
}

// Mock listener for testing
type mockListener struct {
	acceptChan chan net.Conn
	closed     bool
}

func newMockListener() *mockListener {
	return &mockListener{
		acceptChan: make(chan net.Conn),
		closed:     false,
	}
}

func (m *mockListener) Accept() (net.Conn, error) {
	conn, ok := <-m.acceptChan
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

func (m *mockListener) Close() error {
	if !m.closed {
		close(m.acceptChan)
		m.closed = true
	}
	return nil
}

func (m *mockListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
}

var netListen = net.Listen

// Test timeout handling
func TestAddNewClientTimeout(t *testing.T) {
	mock := &mockConn{}

	done := make(chan bool)
	go func() {
		AddNewClient(mock)
		done <- true
	}()

	// Wait for timeout or completion
	select {
	case <-done:
		if !mock.closed {
			t.Error("Connection should be closed after timeout")
		}
	case <-time.After(65 * time.Second):
		t.Error("Function did not timeout as expected")
	}
}

// Test concurrent access
func TestAddNewClientConcurrent(t *testing.T) {
	const numClients = 10
	done := make(chan bool, numClients)

	for i := 0; i < numClients; i++ {
		go func(id int) {
			mock := newMockConn(fmt.Sprintf("user%d\n", id))
			AddNewClient(mock)
			done <- true
		}(i)
	}

	// Wait for all clients to complete
	for i := 0; i < numClients; i++ {
		select {
		case <-done:
			// Client completed successfully
		case <-time.After(5 * time.Second):
			t.Error("Timeout waiting for client to complete")
		}
	}

	// Verify no duplicate names were added
	mClients.Lock()
	names := make(map[string]bool)
	for _, client := range clients {
		if names[client.name] {
			t.Error("Duplicate name found:", client.name)
		}
		names[client.name] = true
	}
	mClients.Unlock()
}

// Test invalid input handling
func TestAddNewClientInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"Very long name", strings.Repeat("a", 1025) + "\n", true},
		{"Empty name field", "\n", true},
		{"Special characters", "user@#$%\n", false},
		{"Unicode characters", "用户名\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockConn(tt.input)
			AddNewClient(mock)

			output := getWrittenData(mock)
			if tt.wantError && !strings.Contains(output, "error") {
				t.Error("Expected error message not found in output")
			}
		})
	}
}

func TestServer(t *testing.T) {
	tests := []struct {
		name string
		port string
		// mockListener *mockListener
		wantError bool
	}{
		{"invalid port", ":invalid", true},
		{"valid port", ":8080", false},
	}

	originalNetListen := netListen
	defer func() { netListen = originalNetListen }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			netListen = func(network, addr string) (net.Listener, error) {
				if tt.wantError {
					return nil, net.ErrClosed
				}
				return newMockListener(), nil
			}

			serverStarted := make(chan struct{})

			go func() {
				close(serverStarted)
				Server(tt.port)
			}()

			select {
			case <-serverStarted:
			case <-time.After(time.Second):

			}
		})
	}
}
