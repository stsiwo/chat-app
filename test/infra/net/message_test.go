package net

import (
  "testing"
  "github.com/stsiwo/chat-app/domain/user"
  cnet "github.com/stsiwo/chat-app/infra/net"
  "github.com/stsiwo/chat-app/mocks"
	"github.com/gobwas/ws"
  "net"
  "encoding/json"
	"github.com/stretchr/testify/assert"
  "time"
)

func TestMessageEncodingDecoding(t *testing.T) {

  _, dummyConn := net.Pipe()

  dummySender := cnet.NewClient(
    dummyConn,
    user.NewAdminUser("sample-admin-user"),
    &mocks.IPool{},
    &mocks.IPool{},
    &mocks.IWsutilHandler{},
  )

  dummyReceiver := cnet.NewClient(
    dummyConn,
    user.NewGuestUser("sample-guest-user"),
    &mocks.IPool{},
    &mocks.IPool{},
    &mocks.IWsutilHandler{},
  )

  sutMessage := cnet.NewMessage(
    dummySender,
    dummyReceiver,
    "sample-content",
    cnet.Text,
  )

  sutMessage.SetOpcode(byte(ws.OpPong))

  jsoned, _ := json.Marshal(sutMessage)
  var newDummyMessage *cnet.Message
  json.Unmarshal(jsoned, &newDummyMessage)

  assert.Equal(t, sutMessage.Id(), newDummyMessage.Id())
  assert.Equal(t, sutMessage.Content(), newDummyMessage.Content())
  assert.Equal(t, sutMessage.Date().Format(time.RFC3339), newDummyMessage.Date().Format(time.RFC3339))
  assert.Equal(t, sutMessage.Opcode(), newDummyMessage.Opcode())
}
