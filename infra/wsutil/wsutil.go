package wsutil

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/gobwas/ws"
  "io"
)

type IWsutilHandler interface {
  ReadClientData(rw io.ReadWriter) ([]byte, ws.OpCode, error)
  WriteServerMessage(w io.Writer, op ws.OpCode, p []byte) error
}

type Handler struct {
}

func (h Handler) ReadClientData(rw io.ReadWriter) ([]byte, ws.OpCode, error) {
  return wsutil.ReadClientData(rw)
}

func (h Handler) WriteServerMessage(w io.Writer, op ws.OpCode, p []byte) error {
  return wsutil.WriteServerMessage(w, op, p)
}
