package functest

import (
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stsiwo/chat-app/infra/net"
	"github.com/stsiwo/chat-app/ui"
	"net/http/httptest"
	"testing"
	//"net/http"
	"fmt"
	"log"
	//"io/ioutil"
	"context"
	"github.com/gobwas/ws"
	"strings"
	//"runtime"
	"encoding/json"
	"time"
)

func TestFuncUpgradeRequestShouldReturnConn(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	urlqs := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=1&id=" + uuid.New().String()
	fmt.Printf("ws url: %v\n", urlqs)

	// get ws connection
	conn, _, _, err := ws.Dial(context.Background(), urlqs)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t, conn != nil)
	// make sure all necessary GRs are running
	// pool * 2 + 1(num of client) * (1 (read GR) + 1 (write GR)) = 4
	// need profiling tool
	//assert.Equal(t, 5, runtime.NumGoroutine())
}

func TestFuncShouldRegisterAdminClientToAdminPool(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	urlqs := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=2&id=" + uuid.New().String()
	fmt.Printf("ws url: %v\n", urlqs)

	// get ws connection
	conn, _, _, err := ws.Dial(context.Background(), urlqs)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t, conn != nil)
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 1, adminPool.Size())
}

func TestFuncShouldRegisterMemberClientToUserPool(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	urlqs := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=1&id=" + uuid.New().String()
	fmt.Printf("ws url: %v\n", urlqs)

	// get ws connection
	conn, _, _, err := ws.Dial(context.Background(), urlqs)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t, conn != nil)
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 1, userPool.Size())
}

func TestFuncShouldRegisterGuestClientToUserPool(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	urlqs := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=0&id=" + uuid.New().String()
	fmt.Printf("ws url: %v\n", urlqs)

	// get ws connection
	conn, _, _, err := ws.Dial(context.Background(), urlqs)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t, conn != nil)
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 1, userPool.Size())
}

func TestFuncShouldSendMessageFromUserToAllAdminClient(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	guestId := uuid.New().String()
	guestUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=0&id=" + guestId
	fmt.Printf("ws url: %v\n", guestUrl)
	adminId := uuid.New().String()
	adminUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=2&id=" + adminId
	fmt.Printf("ws url: %v\n", adminUrl)
	admin1Id := uuid.New().String()
	admin1Url := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=2&id=" + admin1Id
	fmt.Printf("ws url: %v\n", adminUrl)

	// get ws connection
	guestConn, _, _, err := ws.Dial(context.Background(), guestUrl)
	adminConn, _, _, err := ws.Dial(context.Background(), adminUrl)
	admin1Conn, _, _, err := ws.Dial(context.Background(), admin1Url)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	// guestClient
	guestClient := userPool.Find(guestId)

	// dummy message
	dummyMessage := net.NewMessage(guestClient, nil, "hello, world", net.Text)
	jsonedMessage, err := json.Marshal(dummyMessage)
	if err != nil {
		log.Fatalf("err during encoding dummy message: %v \n", err)
	}

	// send message from guest to admin
	wsutil.WriteClientMessage(guestConn, 0x1, jsonedMessage)

  // receiving message from admin connection
	revRawMsg, _ /*op*/, err := wsutil.ReadServerData(adminConn)
	revRawMsg1, _ /*op*/, err := wsutil.ReadServerData(admin1Conn)

  // decoding to Message instnace
	var message net.Message
	err = json.Unmarshal(revRawMsg, &message)
	if err != nil {
		log.Fatalf("json decoding error when reading dummy received message: %v \n", err)
	}

	var message1 net.Message
	err = json.Unmarshal(revRawMsg1, &message1)
	if err != nil {
		log.Fatalf("json decoding error when reading dummy received message: %v \n", err)
	}

  // check original message matches with received message
  assert.Equal(t, dummyMessage.Id(), message.Id())
  assert.Equal(t, dummyMessage.Content(), message.Content())

  assert.Equal(t, dummyMessage.Id(), message1.Id())
  assert.Equal(t, dummyMessage.Content(), message1.Content())
}

func TestFuncShouldSendMessageFromAdminToSpecificUserClient(t *testing.T) {

	// prepare & run pools
	adminPool := net.NewPool()
	userPool := net.NewPool()
	go adminPool.Run()
	go userPool.Run()

	// setup controller
	wsc := ui.NewController(adminPool, userPool)

	// run test server
	ts := httptest.NewServer(wsc)
	defer ts.Close()

	// prep server url with query string
	guest1Id := uuid.New().String()
	guest1Url := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=0&id=" + guest1Id
	fmt.Printf("ws url: %v\n", guest1Url)
	guest2Id := uuid.New().String()
	guest2Url := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=0&id=" + guest2Id
	fmt.Printf("ws url: %v\n", guest2Url)
	adminId := uuid.New().String()
	adminUrl := "ws" + strings.TrimPrefix(ts.URL, "http") + "?role=2&id=" + adminId
	fmt.Printf("ws url: %v\n", adminUrl)

	// get ws connection
	guest1Conn, _, _, err := ws.Dial(context.Background(), guest1Url)
	adminConn, _, _, err := ws.Dial(context.Background(), adminUrl)
	// adding below causes "EOF" error at Read()
	//defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	// guestClient
	guest1Client := userPool.Find(guest1Id)
	adminClient := adminPool.Find(adminId)

	// dummy admin message
	dummyMessage := net.NewMessage(adminClient, guest1Client, "hello, world", net.Text)
	jsonedMessage, err := json.Marshal(dummyMessage)
	if err != nil {
		log.Fatalf("err during encoding dummy message: %v \n", err)
	}

	// send message from admin to guest1
	wsutil.WriteClientMessage(adminConn, 0x1, jsonedMessage)

  // receiving message from guest1 conn 
	revRawMsg, _ /*op*/, err := wsutil.ReadServerData(guest1Conn)

  // decoding to Message instnace
	var message net.Message
	err = json.Unmarshal(revRawMsg, &message)
	if err != nil {
		log.Fatalf("json decoding error when reading dummy received message: %v \n", err)
	}

  // check original message matches with received message
  assert.Equal(t, dummyMessage.Id(), message.Id())
  assert.Equal(t, dummyMessage.Content(), message.Content())
  assert.True(t, false)
}
