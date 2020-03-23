package net

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/stsiwo/chat-app/domain/user"
	"log"
	"net"
)

type IClient interface {
  Read()
  Write()
  Id() string
  Send(message *Message)
  Receive() *Message
}

type Client struct {
	id string

	conn net.Conn

	user *user.User

	send chan *Message

	adminPool IPool

	userPool IPool
}

func NewClient(conn net.Conn, user *user.User, adminPool IPool, userPool IPool) *Client {
	return &Client{
		id:   uuid.New().String(),
		conn: conn,
		user: user,
		// need to set buffer size >= 2
		// since 'broadcasting' inside for loop at pool.go will be blocked when
		// sending message to each client and can't go through each iteration
		// currently set to 2 but need to increase based on how many messages for a single client can keep those message
		send:      make(chan *Message, 2),
		adminPool: adminPool,
		userPool:  userPool,
	}
}

func (c *Client) Id() string {
  return c.id
}

func (c *Client) Send(m *Message) {
  c.send <- m
}

func (c *Client) Receive() *Message {
  return <-c.send
}

// GR
func (c *Client) Read() {
	defer c.conn.Close()

	for {
    // wait for new message; block until new message arrives 
		rowMsg, opcode, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			log.Fatalf("reading client data error of wsutil package: %v\n", err)
		}

    // concert []byte message to json and also convert it to Message object
		var message *Message
		err = json.Unmarshal(rowMsg, message)
		if err != nil {
			log.Fatalf("json decoding error when reading message: %v \n", err)
		}

		// set opcode to the message
		message.setOpcode(byte(opcode))

    // send message to proper destination via channel
		if c.user.Role() == user.Admin {
      // if current client is admin; receive new message from admin browser, send message to the specified user client (unicast)
			c.userPool.Unicast(message)
		} else {
      // if current client is user; receive new message from user browser, send message to all admin users (broadcast)
			c.adminPool.Broadcast(message)
		}
	}
}

// GR
func (c *Client) Write() {
	defer c.conn.Close()

	for {

		select {
		case m := <-c.send:
			// extract opcode from received message
			opcode := ws.OpCode(m.Opcode())

      // convert message struct to json
      jsonMsg, err := json.Marshal(m)
      if err != nil {
        log.Fatalf("json encoding error when writing message: %v \n", err)
      }

      // convert json message to []byte
			err = wsutil.WriteServerMessage(c.conn, opcode, []byte(jsonMsg))
			if err != nil {
        log.Fatalf("error during writing message to client: %v \n", err)
			}
		}
	}
}
