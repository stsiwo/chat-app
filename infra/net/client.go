package net

import (
  "net"
  "github.com/stsiwo/chat-app/domain/user"
	"github.com/google/uuid"
)

type Client struct {

  id string

  conn net.Conn

  user *user.User

  send chan []byte
}

func NewClient(conn net.Conn, user *user.User) *Client {
  return &Client{
    id: uuid.New().String(),
    conn: conn,
    user: user,
    // need to set buffer size >= 2
    // since 'broadcasting' inside for loop at pool.go will be blocked when 
    // sending message to each client and can't go through each iteration
    // currently set to 2 but need to increase based on how many messages for a single client can keep those message
    send: make(chan []byte, 2),
  }
}
