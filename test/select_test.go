package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestSelect(t *testing.T) {

	var ws sync.WaitGroup
	c1 := make(chan int)
	c2 := make(chan int)

	ws.Add(2)
	go func() {
		for {
			fmt.Println("inside loop")
			select { // select block one case has met
			case i := <-c1:
				fmt.Println("c1 channel with value: " + strconv.Itoa(i))
				ws.Done()
			case i := <-c2:
				fmt.Println("c2 channel with value: " + strconv.Itoa(i))
				ws.Done()
				//default: // run if no other case is ready
				//fmt.Println("default")
				//ws.Done()
			}
		}
	}()

	c2 <- 1
	c1 <- 1

	ws.Wait()

	assert.True(t, true)
}
