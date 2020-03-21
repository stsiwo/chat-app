package net

import (
	"encoding/json"
	//"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/stsiwo/chat-app/domain/user"
	"log"
	"net"
)

type Client struct {
	id string

	conn net.Conn

	user *user.User

	send chan *Message

	adminPool *Pool

	userPool *Pool
}

func NewClient(conn net.Conn, user *user.User, adminPool *Pool, userPool *Pool) *Client {
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

// GR
func (c *Client) read() {
	defer c.conn.Close()

	//var (
	//	r       = wsutil.NewReader(c.conn, ws.StateServerSide)
	//	w       = wsutil.NewWriter(c.conn, ws.StateServerSide, ws.OpText)
	//	decoder = json.NewDecoder(r)
	//	encoder = json.NewEncoder(w)
	//)

	for {
		rowMsg, _, err := wsutil.ReadClientData(c.conn)
		if err != nil {
      log.Fatalf("reading client data error of wsutil package: %v\n", err)
		}

		var message *Message
		err = json.Unmarshal(rowMsg, message)
		if err != nil {
      log.Fatalf("json decoding error when reading message: %v \n", err)
		}

		if c.user.Role() == user.Admin {
			c.userPool.unicast <-message
		} else {
			c.adminPool.broadcast <- message
		}
	}
}

// GR
func (c *Client) write() {
}
