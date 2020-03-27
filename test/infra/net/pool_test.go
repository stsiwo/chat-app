package net

import (
	"github.com/google/uuid"
	"github.com/stsiwo/chat-app/domain/user"
	cnet "github.com/stsiwo/chat-app/infra/net"
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
	var pool cnet.IPool = cnet.NewPool()

	var dummyClient cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyConn,
		user.NewGuestUser(uuid.New().String(), "test-user-name"),
		pool,
		nil,
		nil,
	)

	go pool.Run()

	ws.Add(1)
	go func() {
		defer ws.Done()
		pool.Register(dummyClient)
	}()

	ws.Wait()

	assert.Equal(t, 1, pool.Size())
}

func TestPoolRegisterShouldStoreMultipleClients(t *testing.T) {

	var ws sync.WaitGroup
	var dummyClientList [100]cnet.IClient
	var pool cnet.IPool = cnet.NewPool()

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = cnet.NewClient(
			uuid.New().String(),
			dummyConn,
			user.NewGuestUser(uuid.New().String(), "test-user-name"+strconv.Itoa(i)),
			pool,
			nil,
			nil,
		)
	}

	go pool.Run()

	for _, c := range dummyClientList {

		ws.Add(1)
		// don't foreget set parameters 'c' otherwise wierd error
		go func(c cnet.IClient) {
			defer ws.Done()
			pool.Register(c)
		}(c)
	}

	ws.Wait()

	// ?? length = 100 but result shows only 99
	// please fix this once you figure out
	// fixed but i don't know why
	// use another GR like above. it solves this error
	// ?? still produce this error sometimes
	assert.Equal(t, len(dummyClientList), pool.Size())
}

func TestPoolFindShouldGetSpecifiedClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	var pool cnet.IPool = cnet.NewPool()

	var dummyClient cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyConn,
		user.NewGuestUser(uuid.New().String(), "test-user-name"),
		pool,
		nil,
		nil,
	)

	go pool.Run()

	ws.Add(1)
	var receivedClient cnet.IClient
	go func() {
		defer ws.Done()
		pool.Register(dummyClient)
		receivedClient = pool.Find(dummyClient.Id())
		_ = receivedClient // skip 'declare but not used' compile error
	}()

	ws.Wait()

	assert.Equal(t, dummyClient.Id(), receivedClient.Id())
}

func TestPoolUnregisterShouldRemoveSpecifiedClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	var pool cnet.IPool = cnet.NewPool()

	var dummyClient cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyConn,
		user.NewGuestUser(uuid.New().String(), "test-user-name"),
		pool,
		nil,
		nil,
	)

	go pool.Run()

	ws.Add(1)
	go func(pool cnet.IPool) {
		defer ws.Done()
		pool.Register(dummyClient)
		receivedClient := pool.Find(dummyClient.Id())
		pool.Unregister(receivedClient)
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
	assert.Equal(t, 0, pool.Size())
}

func TestPoolBroadcastShouldDeliverMessageToPoolWithSingleClient(t *testing.T) {
	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	var pool cnet.IPool = cnet.NewPool()

	var dummyClient cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyConn,
		user.NewGuestUser(uuid.New().String(), "test-user-name"),
		pool,
		nil,
		nil,
	)

	dummyMessage := cnet.NewMessage(
		dummyClient,
		nil,
		"sample-message-content",
		cnet.Text,
	)

	go pool.Run()
	ws.Add(1)
	go func(pool cnet.IPool) {
		defer ws.Done()
		pool.Register(dummyClient)
		pool.Broadcast(dummyMessage)
	}(pool)

	ws.Wait()

	expectedMessage := dummyClient.Receive()

	assert.Equal(t, expectedMessage, dummyMessage)

}

func TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient(t *testing.T) {
	fmt.Println("\nstart TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient")

	var ws sync.WaitGroup
	var dummyClientList [10]cnet.IClient
	var pool cnet.IPool = cnet.NewPool()

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = cnet.NewClient(
			uuid.New().String(),
			dummyConn,
			user.NewGuestUser(uuid.New().String(), "test-user-name"+strconv.Itoa(i)),
			nil,
			pool,
			nil,
		)
	}

	go pool.Run()

	for _, c := range dummyClientList {
		ws.Add(1)
		// don't foreget set parameters 'c' otherwise wierd error
		go func(c cnet.IClient) {
			defer ws.Done()
			pool.Register(c)
		}(c)
	}

	ws.Wait()

	dummyMessage := cnet.NewMessage(
		dummyClientList[0],
		nil,
		"sample-message-content",
		cnet.Text,
	)
	ws.Add(1)
	go func() {
		defer ws.Done()
		pool.Broadcast(dummyMessage)
	}()

	ws.Wait()

	for _, c := range dummyClientList {
		fmt.Printf("client: %v \n", c)
		expectedMessage := c.Receive()
		assert.Equal(t, expectedMessage, dummyMessage)
	}
}

func TestPoolUnicastShouldDeliverMessageToSpecificClientInPool(t *testing.T) {
	fmt.Println("\nstart TestPoolUnicastShouldDeliverMessageToSpecificClientInPool")

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()
	_, dummyReceiverConn := net.Pipe()
	var pool cnet.IPool = cnet.NewPool()

	var dummyClient cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyConn,
		user.NewGuestUser(uuid.New().String(), "test-user-name"),
		pool,
		nil,
		nil,
	)

	var dummyReceiver cnet.IClient = cnet.NewClient(
		uuid.New().String(),
		dummyReceiverConn,
		user.NewGuestUser(uuid.New().String(), "test-receiver-name"),
		pool,
		nil,
		nil,
	)

	dummyMessage := cnet.NewMessage(
		dummyClient,
		dummyReceiver,
		"sample-message-content",
		cnet.Text,
	)

	go pool.Run()

	ws.Add(1)
	go func(pool cnet.IPool) {
		defer ws.Done()
		pool.Register(dummyClient)
		pool.Register(dummyReceiver)
		pool.Unicast(dummyMessage)
	}(pool)

	ws.Wait()

	expectedMessage := dummyReceiver.Receive()

	assert.Equal(t, expectedMessage, dummyMessage)
}
