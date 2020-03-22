package net

import (
	"github.com/google/uuid"
	"github.com/stsiwo/chat-app/domain/user"
	//"github.com/stsiwo/chat-app/domain/main"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"sync"
	"testing"
	//"time"
	"runtime"
)

func TestPoolRegisterShouldStoreNewClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	pool := NewPool()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("test-user-name"),
		pool,
		nil,
	)

	go pool.Run()

	ws.Add(1)
	go func() {
		defer ws.Done()
		pool.register <- dummyClient
	}()

	ws.Wait()

	assert.Equal(t, 1, len(pool.pool))
}

func TestPoolRegisterShouldStoreMultipleClients(t *testing.T) {

	var ws sync.WaitGroup
	var dummyClientList [100]*Client
	pool := NewPool()

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = NewClient(
			dummyConn,
			user.NewGuestUser("test-user-name"+strconv.Itoa(i)),
			pool,
			nil,
		)
	}

	go pool.Run()

	for _, c := range dummyClientList {

		ws.Add(1)
		// don't foreget set parameters 'c' otherwise wierd error
		go func(c *Client) {
			defer ws.Done()
			pool.register <- c
		}(c)
	}

	ws.Wait()

	// ?? length = 100 but result shows only 99
	// please fix this once you figure out
	// fixed but i don't know why
	// use another GR like above. it solves this error
	// ?? still produce this error sometimes
	assert.Equal(t, len(dummyClientList), len(pool.pool))
}

func TestPoolFindShouldGetSpecifiedClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	pool := NewPool()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("test-user-name"),
    pool,
    nil,
	)

	go pool.Run()

	ws.Add(1)
	var receivedClient *Client
	go func() {
		defer ws.Done()
		pool.register <- dummyClient
		receivedClient = pool.find(dummyClient.id)
		_ = receivedClient // skip 'declare but not used' compile error
	}()

	ws.Wait()

	assert.Equal(t, dummyClient.id, receivedClient.id)
}

func TestPoolUnregisterShouldRemoveSpecifiedClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	pool := NewPool()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("test-user-name"),
    pool,
    nil,
	)

	go pool.Run()

	ws.Add(1)
	go func(pool *Pool) {
		defer ws.Done()
		pool.register <- dummyClient
		receivedClient := pool.find(dummyClient.id)
		pool.unregister <- receivedClient
	}(pool)

	ws.Wait()
	runtime.GC() // make sure to delete item is collected

	/**
	 * issue: can't delete item inside pool.pool (map) item
	 * details:
	 *  - can delete item inside 'run()' at pool.go file, but pool object at this test still holding deleted item.
	 * assumptions:
	 *  - two pool.pool (map) objects at this testing and the one at 'pool.go' refer to different memory so deleting one of them does not affect to the main (this testing object)
	 **/
	/**
	 * solved above problem!!!!
	 * - add runtime.GC()
	 **/

	/**
	 * another issue: sometimes, running this test produce below error:
	 * panic: runtime error: invalid memory address or nil pointer dereference
	 **/
	assert.Equal(t, 0, len(pool.pool))
}

func TestPoolBroadcastShouldDeliverMessageToPoolWithSingleClient(t *testing.T) {
	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	pool := NewPool()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("test-user-name"),
    pool,
    nil,
	)

	dummyMessage := NewMessage(
    dummyClient,
    nil,
    "sample-message-content",
  )


	go pool.Run()
	ws.Add(1)
	go func(pool *Pool) {
		defer ws.Done()
		pool.register <- dummyClient
		pool.broadcast <- dummyMessage
	}(pool)

	ws.Wait()

	expectedMessage := <-dummyClient.send

	assert.Equal(t, expectedMessage, dummyMessage)

}

func TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient(t *testing.T) {
	fmt.Println("\nstart TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient")

	var ws sync.WaitGroup
	var dummyClientList [10]*Client
	pool := NewPool()

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = NewClient(
			dummyConn,
			user.NewGuestUser("test-user-name"+strconv.Itoa(i)),
      nil,
      pool,
		)
	}

	go pool.Run()

	for _, c := range dummyClientList {
		ws.Add(1)
		// don't foreget set parameters 'c' otherwise wierd error
		go func(c *Client) {
			defer ws.Done()
			pool.register <- c
		}(c)
	}

	ws.Wait()

	dummyMessage := NewMessage(
    dummyClientList[0],
    nil,
    "sample-message-content",
  )
	ws.Add(1)
	go func() {
		defer ws.Done()
		pool.broadcast <- dummyMessage
	}()

	ws.Wait()

	for _, c := range dummyClientList {
		fmt.Printf("client: %v \n", c)
		expectedMessage := <-c.send
		assert.Equal(t, expectedMessage, dummyMessage)
	}
}

func TestPoolUnicastShouldDeliverMessageToSpecificClientInPool(t *testing.T) {
	fmt.Println("\nstart TestPoolUnicastShouldDeliverMessageToSpecificClientInPool")

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	_, dummyReceiverConn := net.Pipe()
	pool := NewPool()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("test-user-name"),
    pool,
    nil,
	)

  dummyReceiver := NewClient(
		dummyReceiverConn,
		user.NewGuestUser("test-receiver-name"),
    pool,
    nil,
  )

	dummyMessage := NewMessage(
    dummyClient,
    dummyReceiver,
    "sample-message-content",
  )

	go pool.Run()

	ws.Add(1)
	go func(pool *Pool) {
		defer ws.Done()
		pool.register <- dummyClient
		pool.register <- dummyReceiver
		pool.unicast <- dummyMessage
	}(pool)

	ws.Wait()

	expectedMessage := <-dummyReceiver.send

	assert.Equal(t, expectedMessage, dummyMessage)
}
