package net

import (
	"github.com/google/uuid"
  "time"
)

type Message struct {

  id string

  sender IClient

  receiver IClient

  // any type => 'interface{}' and use type switch for handling each type
  content interface{}

  date time.Time

  opcode byte

}

func NewMessage(sender IClient, receiver IClient, content interface{}) *Message {
  return &Message{
    id: uuid.New().String(),
    sender: sender,
    receiver: receiver,
    content: content,
    date: time.Now(),
    opcode: byte(0),
  }
}

func (m *Message) Id() string {
  return m.id
}

func (m *Message) Sender() IClient {
  return m.sender
}

func (m *Message) Receiver() IClient {
  return m.receiver
}

func (m *Message) Content() interface{} {
  return m.content
}

func (m *Message) Date() time.Time {
  return m.date
}

func (m *Message) Opcode() byte {
  return m.opcode
}

func (m *Message) setOpcode(opcode byte) {
  m.opcode = opcode
}


