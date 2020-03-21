package test

import (
  "testing"
  //"sync"
  "fmt"
  "github.com/stretchr/testify/assert"
  //"strconv"
)


//func TestChannel(t *testing.T) {
//
//  fmt.Println("start channel test")
//
//  channel := make(chan int, 3)
//
//  go func() {
//    channel <- 1
//    channel <- 2
//    channel <- 3
//    fmt.Println("before 4")
//    channel <- 4
//    fmt.Println("after 4")
//    close(channel)
//  }()
//
//  go func() {
//    data := <-channel
//    fmt.Println("data: " + strconv.Itoa(data))
//  }()
//
//  fmt.Println("end channel test")
//
//  assert.True(t, true)
//
//}
//
//type Client struct {
//
//  value string
//
//}
//
//func TestChannelWithPointer(t *testing.T) {
//
//  var ws sync.WaitGroup
//  channel := make(chan *Client)
//
//  ws.Add(1)
//  go func() {
//    for {
//      select {
//      case c := <-channel:
//        c.value = "updated"
//        ws.Done()
//      }
//    }
//  }()
//
//  dummyClient := &Client{value: "satoshi"}
//
//  channel <-dummyClient
//
//  ws.Wait()
//
//  fmt.Println(dummyClient.value)
//
//  assert.True(t, true)
//}

func TestChannelWaitUntilReceiveData(t *testing.T) {

  //var ws sync.WaitGroup
  channel := make(chan string, 3)

  //ws.Add(1)
  go func(channel chan string) {

    data1 := <-channel
    data2 := <-channel
    data3 := <-channel
    fmt.Printf("received data from channel: %v \n", data1)
    fmt.Printf("received data from channel: %v \n", data2)
    fmt.Printf("received data from channel: %v \n", data3)

  }(channel)
  channel <-"my-data"
  channel <-"my-data"
  channel <-"my-data"

  assert.True(t, false)
}
