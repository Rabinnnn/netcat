package utils

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

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

// Test invalid input handling
func TestAddNewClientInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "Very long name",
			input:     strings.Repeat("a", 1025) + "\n",
			wantError: true,
		},
		{
			name:      "Empty name field",
			input:     "\n",
			wantError: true,
		},
		{
			name:      "Special characters",
			input:     "user@#$%\n",
			wantError: false,
		},
		{
			name:      "Unicode characters",
			input:     "用户名\n",
			wantError: false,
		},
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
