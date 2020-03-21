package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
  "runtime"
  "fmt"
)

func TestMapShouldRemoveItem(t *testing.T) {

  myMap := make(map[string]int)

  myMap["1"] = 1

  delete(myMap, "1")


	assert.Equal(t, 0, len(myMap))
}

type SamplePool struct {

  pool map[string]string

}

func NewSamplePool() SamplePool {
  return SamplePool{
    pool: make(map[string]string),
  }
}

func (s *SamplePool) remove(key string) {
  delete(s.pool, key)
}

func TestMapWithStruct(t *testing.T) {

  PrintMemUsage()
  sp := NewSamplePool()
  PrintMemUsage()

  sp.pool["1"] = "1"

  sp.remove("1")


  runtime.GC()
  PrintMemUsage()
  assert.Equal(t, 1, len(sp.pool))
  PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v B", m.Alloc)
        fmt.Printf("\tTotalAlloc = %v B", m.TotalAlloc)
        fmt.Printf("\tSys = %v B", m.Sys)
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

//func TestMapWithChannel(t *testing.T) {
//
//  var ws sync.
//  sutPool := TestPool{
//    pool: make(map[string]Test),
//    unregister: make(chan Test),
//  }
//
//  sutPool.pool["1"] = Test{value: "satoshi"}
//
//  go func() {
//    for {
//      select {
//      case t := <-sutPool.unregister:
//        delete(sutPool.pool, "1")
//        _ = t
//      }
//    }
//  }()
//}
