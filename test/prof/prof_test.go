package prof

import (
  "sync"
  "time"
  "fmt"
  "testing"
)


func TestProf(t *testing.T) {
    fmt.Println("hello world")
    var wg sync.WaitGroup
    wg.Add(1)
    go leakyFunction(wg)
    wg.Wait()
}

func leakyFunction(wg sync.WaitGroup) {
    defer wg.Done()
    s := make([]string, 3)
    for i:= 0; i < 10000000; i++{
        s = append(s, "magical pandas")
        if (i % 100000) == 0 {
            time.Sleep(500 * time.Millisecond)
        }
    }
}

