package net

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type MessageType int8

const (
	Text   MessageType = 0
	Binary MessageType = 1
)

type Message struct {
	id string

	messageType MessageType

	sender IClient

	receiver IClient

	// any type => 'interface{}' and use type switch for handling each type
	content interface{}

	date time.Time

	opcode byte
}

func NewMessage(sender IClient, receiver IClient, content interface{}, messageType MessageType) *Message {
	return &Message{
		id:          uuid.New().String(),
		messageType: messageType,
		sender:      sender,
		receiver:    receiver,
		content:     content,
		date:        time.Now(),
		opcode:      byte(0),
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

func (m *Message) SetOpcode(opcode byte) {
	m.opcode = opcode
}

func (m *Message) UnmarshalJSON(rawData []byte) error {
	log.Println("unmarshal json is called inside message struct")
	var objMap map[string]json.RawMessage
	err := json.Unmarshal(rawData, &objMap)
	if err != nil {
    log.Printf("error during converting []bytes to map struct at message struct: %v", err)
		return err
	}
	// need to remove double quote from rawMessage
	m.id = strings.Trim(string(objMap["id"]), "\"")
	messageTypeNum, err := strconv.Atoi(string(objMap["messageType"]))
	if err != nil {
    log.Printf("error during converting messageType at message struct: %v", err)
		return err
	}
	m.messageType = MessageType(messageTypeNum)
	var tempSender *Client
	json.Unmarshal(objMap["sender"], &tempSender)
	m.sender = tempSender
	var tempReceiver *Client
	json.Unmarshal(objMap["receiver"], &tempReceiver)
	m.receiver = tempReceiver
	if m.messageType == Binary {
		m.content = objMap["content"]
	} else if m.messageType == Text {
		m.content = strings.Trim(string(objMap["content"]), "\"")
	}
	m.date, err = time.Parse(time.RFC3339, strings.Trim(string(objMap["date"]), "\""))
	log.Printf("opcode: %v", objMap["opcode"])
	opcodeNum, err := strconv.Atoi(string(objMap["opcode"]))
	if err != nil {
    log.Printf("error during converting opcode at message struct: %v", err)
		return err
	}
	m.opcode = byte(opcodeNum)
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":          m.id,
		"messageType": m.messageType,
		"sender":      m.sender,
		"receiver":    m.receiver,
		"content":     m.content,
		"date":        m.date.Format(time.RFC3339),
		"opcode":      m.opcode,
	})
}

// return string type message content if message type is Text
// otherwise return empty string
func (m *Message) GetTextConent() (string, error) {
	if m.messageType == Text {
		return m.content.(string), nil
	}
	return "", errors.New("you got wrong type of message: expected Text")
}

// return []byte type message content if message type is Binary
// otherwise return 0 byte
func (m *Message) GetBinaryConent() ([]byte, error) {
	if m.messageType == Binary {
		return m.content.([]byte), nil
	}
	return []byte(""), errors.New("you got wrong type of message: expected Binary")
}
