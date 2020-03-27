package net

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/stsiwo/chat-app/domain/user"
	"github.com/stsiwo/chat-app/infra/wsutil"
	"log"
	"net"
	"strings"
)

type IClient interface {
	Read()
	Write()
	Id() string
	Send(message *Message)
	Receive() *Message
	User() *user.User
}

type Client struct {
	id string

	conn net.Conn

	user *user.User

	send chan *Message

	adminPool IPool

	userPool IPool

	wsutilHandler wsutil.IWsutilHandler
}

func NewClient(id string, conn net.Conn, user *user.User, adminPool IPool, userPool IPool, wsutilHandler wsutil.IWsutilHandler) *Client {
	return &Client{
		id:   id,
		conn: conn,
		user: user,
		// need to set buffer size >= 2
		// since 'broadcasting' inside for loop at pool.go will be blocked when
		// sending message to each client and can't go through each iteration
		// currently set to 2 but need to increase based on how many messages for a single client can keep those message
		send:          make(chan *Message, 2),
		adminPool:     adminPool,
		userPool:      userPool,
		wsutilHandler: wsutilHandler,
	}
}

func (c *Client) Id() string {
	return c.id
}

func (c *Client) User() *user.User {
	return c.user
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
	log.Println("start Read() operation")

	for {
		// wait for new message; block until new message arrives
		rowMsg, opcode, err := c.wsutilHandler.ReadClientData(c.conn)
		if err != nil {
			log.Fatalf("reading client data error of wsutil package: %v\n", err)
		}

		// concert []byte message to json and also convert it to Message object
		// be careful not to do 'var message *Message'
		// this message pointer hold nil (since hasn't yet assiged)
		// so 'var message Message' will create empty instance of Message and you can pass its pointer to json UnMarshal
		var message Message
		err = json.Unmarshal(rowMsg, &message)
		if err != nil {
			log.Fatalf("json decoding error when reading message: %v \n", err)
		}

		// set opcode to the message
		message.SetOpcode(byte(opcode))
		log.Printf("decoded message: %v", message)

		// send message to proper destination via channel
		if c.user.Role() == user.Admin {
			// if current client is admin; receive new message from admin browser, send message to the specified user client (unicast)
			c.userPool.Unicast(&message)
		} else {
			// if current client is user; receive new message from user browser, send message to all admin users (broadcast)
			c.adminPool.Broadcast(&message)
		}
	}
}

// GR
func (c *Client) Write() {
	defer c.conn.Close()
	log.Println("start Write() operation")

	for {

		select {
		case m := <-c.send:
			log.Println("this client received message via send channel at Write()")
			// extract opcode from received message
			opcode := ws.OpCode(m.Opcode())

			// convert message struct to json
			log.Printf("convert message to json bytes: %v", m)
			jsonMsg, err := json.Marshal(m)
			if err != nil {
				log.Fatalf("json encoding error when writing message: %v \n", err)
			}

			// convert json message to []byte
			log.Println("start send jsoned message to this client connection")
			err = c.wsutilHandler.WriteServerMessage(c.conn, opcode, []byte(jsonMsg))
			if err != nil {
				log.Fatalf("error during writing message to client: %v \n", err)
			}
		}
	}
}

func (c *Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":   c.id,
		"user": c.user,
	})
}

func (c *Client) UnmarshalJSON(rawData []byte) error {
	log.Println("unmarshal json is called inside client struct")
	var objMap map[string]json.RawMessage
	err := json.Unmarshal(rawData, &objMap)
	if err != nil {
		log.Fatalf("err during decoding user json data: %v", err)
		return err
	}
	// need to remove double quote from rawMessage
	c.id = strings.Trim(string(objMap["id"]), "\"")
	var tempUser *user.User
	json.Unmarshal(objMap["user"], &tempUser)
	c.user = tempUser
	return nil
}

