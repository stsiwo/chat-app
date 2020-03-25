package ws

import (
	"github.com/gobwas/ws"
  "net/http"
  "net"
  "bufio"
)

type IWsHandler interface {
  UpgradeHTTP(r *http.Request, w http.ResponseWriter) (net.Conn, *bufio.ReadWriter, ws.Handshake, error)
}

type Handler struct {
}

func (h Handler) UpgradeHTTP(r *http.Request, w http.ResponseWriter) (net.Conn, *bufio.ReadWriter, ws.Handshake, error) {
  return ws.UpgradeHTTP(r, w)
}

