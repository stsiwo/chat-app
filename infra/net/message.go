package net

import (
	"github.com/google/uuid"
  "time"
)

type Message struct {

  id string

  sender *Client

  receiver *Client

  // any type => 'interface{}' and use type switch for handling each type
  content interface{}

  date time.Time

}

func NewMessage(sender *Client, receiver *Client, content interface{}) *Message {
  return &Message{
    id: uuid.New().String(),
    sender: sender,
    receiver: receiver,
    content: content,
    date: time.Now(),
  }
}

func (m *Message) Id() string {
  return m.id
}

func (m *Message) Sender() *Client {
  return m.sender
}

func (m *Message) Receiver() *Client {
  return m.receiver
}

func (m *Message) Content() interface{} {
  return m.content
}

func (m *Message) Date() time.Time {
  return m.date
}


