package net

import (
	"github.com/google/uuid"
	"github.com/stsiwo/chat-app/domain/user"
	//"github.com/stsiwo/chat-app/domain/main"
	"net"
	"strconv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	//"time"
	"runtime"
)

func TestPoolRegisterShouldStoreNewClient(t *testing.T) {

	var ws sync.WaitGroup
	_, dummyConn := net.Pipe()

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("tets-user", "test-user-name"),
	)

	pool := newPool()

	go pool.run()

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

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = NewClient(
			dummyConn,
			user.NewGuestUser(uuid.New().String(), "test-user-name"+strconv.Itoa(i)),
		)
	}

	pool := newPool()

	go pool.run()

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

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("tets-user", "test-user-name"),
	)

	pool := newPool()

	go pool.run()

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

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("tets-user", "test-user-name"),
	)

	pool := newPool()

	go pool.run()

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

	dummyClient := NewClient(
		dummyConn,
		user.NewGuestUser("tets-user", "test-user-name"),
	)

	dummyMessage := "sample-message"

	pool := newPool()

	go pool.run()
	ws.Add(1)
	go func(pool *Pool) {
		defer ws.Done()
		pool.register <- dummyClient
		pool.broadcast <- []byte(dummyMessage)
	}(pool)

	ws.Wait()

	expectedMessage := <-dummyClient.send

	assert.Equal(t, string(expectedMessage), dummyMessage)

}

func TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient(t *testing.T) {
  fmt.Println("\nstart TestPoolBroadcastShouldDeliverMessageToPoolWithMultipleClient")

	var ws sync.WaitGroup
	var dummyClientList [10]*Client

	for i := range dummyClientList {
		_, dummyConn := net.Pipe()
		dummyClientList[i] = NewClient(
			dummyConn,
			user.NewGuestUser(uuid.New().String(), "test-user-name"+strconv.Itoa(i)),
		)
	}

	pool := newPool()

	go pool.run()

	for _, c := range dummyClientList {
		ws.Add(1)
		// don't foreget set parameters 'c' otherwise wierd error
		go func(c *Client) {
			defer ws.Done()
			pool.register <- c
		}(c)
	}

	ws.Wait()

	dummyMessage := "sample-message"
  ws.Add(1)
	go func() {
		defer ws.Done()
		pool.broadcast <- []byte(dummyMessage)
	}()

  ws.Wait()

	for _, c := range dummyClientList {
    fmt.Printf("client: %v \n", c)
    expectedMessage := <-c.send
		assert.Equal(t, string(expectedMessage), dummyMessage)
	}
}

//func TestConnectionPoolShouldHoldSingleConnection(t *testing.T) {
//  var c = ConnectionPool{Pool: make(map[string]net.Conn)}
//
//  /**
//   * dummy connection
//   * use 'net.Pipe()'
//   * ref: https://stackoverflow.com/questions/30688685/how-does-one-test-net-conn-in-unit-tests-in-golang
//   **/
//  dummConn, _ := net.Pipe()
//  c.Put("id1", dummConn)
//
//  assert.Equal(t, 1, len(c.Pool), "length = 1")
//}
//
//func TestConnectionPoolShouldProvidedRequestedConnByKey(t *testing.T) {
//  var c = ConnectionPool{Pool: make(map[string]net.Conn)}
//
//  /**
//   * dummy connection
//   * use 'net.Pipe()'
//   * ref: https://stackoverflow.com/questions/30688685/how-does-one-test-net-conn-in-unit-tests-in-golang
//   **/
//  dummConn, _ := net.Pipe()
//  c.Put("id1", dummConn)
//
//  receivedConn := c.Find("id1")
//
//  assert.Equal(t, dummConn, receivedConn, "get same dummy connection")
//}
//
//func TestConnectionPoolShouldThreadSafe(t *testing.T) {
//  var c = ConnectionPool{Pool: make(map[string]net.Conn)}
//  // use this to wait until all GRs has finished their role
//  var wg sync.WaitGroup
//
//  // run 10 GRs and each put connection
//  for i := 0; i < 100; i++ {
//    wg.Add(1)
//    go func(i int) {
//      defer wg.Done()
//      dummConn, _ := net.Pipe()
//      c.Put(strconv.Itoa(i), dummConn)
//    }(i) // IFFE; call function at the same time when define
//  }
//
//  // wait until all GRs has finished
//  wg.Wait()
//
//  assert.Equal(t, 100, len(c.Pool), "length = 100")
//}

