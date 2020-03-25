package mocks

import (
  "net"
)

type MockConn struct {
  net.Conn
}

func (m *MockConn) Read(b []byte) (int, error) {
  return 1, nil
}

func (m *MockConn) Write(b []byte) (int, error) {
  return 1, nil
}
