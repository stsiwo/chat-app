package net

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stsiwo/chat-app/domain/user"
	cnet "github.com/stsiwo/chat-app/infra/net"
	"github.com/stsiwo/chat-app/mocks"
	"log"
	"net"
	"testing"
	"time"
)

func TestClientEncodingDecodingClientStruct(t *testing.T) {

	_, dummyConn := net.Pipe()

	dummyClient := cnet.NewClient(
		dummyConn,
		user.NewAdminUser("sample-admin-user"),
		&mocks.IPool{},
		&mocks.IPool{},
		&mocks.IHandler{},
	)

	jsoned, _ := json.Marshal(dummyClient)
	var newDummyClient *cnet.Client
	json.Unmarshal(jsoned, &newDummyClient)

	assert.Equal(t, dummyClient.Id(), newDummyClient.Id())
	assert.Equal(t, dummyClient.User().Id(), newDummyClient.User().Id())
	assert.Equal(t, dummyClient.User().Name(), newDummyClient.User().Name())
}

func TestClientReadShouldSendMessageToSpecificAdminUser(t *testing.T) {

	dummyConn := &mocks.MockConn{}
	dummyAdminUser := user.NewAdminUser("test-admin")
	dummyAdminPool := &mocks.IPool{}
	dummyUserPool := &mocks.IPool{}
	dummyWsHandler := &mocks.IHandler{}
	dummyMessage := cnet.NewMessage(&mocks.IClient{}, &mocks.IClient{}, "sample-message", cnet.Text)
	jsonMsg, err := json.Marshal(dummyMessage)
	if err != nil {
		log.Fatalf("error during json encoding dummy message: %v", err)
	}

	log.Printf("json message: %v \n", jsonMsg)

	dummyWsHandler.On("ReadClientData", dummyConn).Return(jsonMsg, ws.OpPing, nil)
  dummyUserPool.On("Unicast", mock.Anything).Return(nil)

	var sutClient cnet.IClient = cnet.NewClient(
		dummyConn,
		dummyAdminUser,
		dummyAdminPool,
		dummyUserPool,
		dummyWsHandler,
	)

	log.Println("start read go routing")
	go sutClient.Read()

  // need to find proper way to wait until 'Unicast' is called
  // if don't use 'Sleep', the assertion run before the another goroutine finish its job ('Read()')
	time.Sleep(100000 + time.Nanosecond)
	dummyUserPool.AssertCalled(t, "Unicast", mock.Anything)
}

func TestClientReadShouldSendMessageToAllUserClient(t *testing.T) {

	dummyConn := &mocks.MockConn{}
	dummyGuestUser := user.NewGuestUser("test-admin")
	dummyAdminPool := &mocks.IPool{}
	dummyUserPool := &mocks.IPool{}
	dummyWsHandler := &mocks.IHandler{}
	dummyMessage := cnet.NewMessage(&mocks.IClient{}, &mocks.IClient{}, "sample-message", cnet.Text)
	jsonMsg, err := json.Marshal(dummyMessage)
	if err != nil {
		log.Fatalf("error during json encoding dummy message: %v", err)
	}

	dummyWsHandler.On("ReadClientData", dummyConn).Return(jsonMsg, ws.OpPing, nil)
  dummyAdminPool.On("Broadcast", mock.Anything).Return(nil)

	var sutClient cnet.IClient = cnet.NewClient(
		dummyConn,
		dummyGuestUser,
		dummyAdminPool,
		dummyUserPool,
		dummyWsHandler,
	)

	log.Println("start read go routing")
	go sutClient.Read()

  // need to find proper way to wait until 'Unicast' is called
  // if don't use 'Sleep', the assertion run before the another goroutine finish its job ('Read()')
	time.Sleep(100000 + time.Nanosecond)
	dummyAdminPool.AssertCalled(t, "Broadcast", mock.Anything)
}

func TestClientWriteShouldWriteMessageToItsConnection(t *testing.T) {
  // dummy message
  // sut client
	dummyConn := &mocks.MockConn{}
	dummyGuestUser := user.NewGuestUser("test-admin")
	dummyAdminPool := &mocks.IPool{}
	dummyUserPool := &mocks.IPool{}
	dummyWsHandler := &mocks.IHandler{}
	dummyMessage := cnet.NewMessage(&mocks.IClient{}, &mocks.IClient{}, "sample-message", cnet.Text)

	var sutClient cnet.IClient = cnet.NewClient(
		dummyConn,
		dummyGuestUser,
		dummyAdminPool,
		dummyUserPool,
		dummyWsHandler,
	)
  // set behavior of mocked object
	dummyWsHandler.On("WriteServerMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

  // send dummy message to client 'send' channel
  sutClient.Send(dummyMessage)

  // run sut method
  go sutClient.Write()

  // assert IHandler.WriteServerMessage is called
	time.Sleep(1 + time.Millisecond)
	dummyWsHandler.AssertCalled(t, "WriteServerMessage", mock.Anything, mock.Anything, mock.Anything)
}
