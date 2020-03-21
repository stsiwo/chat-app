package main

import (
  "runtime"
  "fmt"
)

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v B", m.Alloc)
        fmt.Printf("\tTotalAlloc = %v B", m.TotalAlloc)
        fmt.Printf("\tSys = %v B", m.Sys)
        fmt.Printf("\tNumGC = %v", m.NumGC)
        fmt.Printf("\tHeapObjects = %v", m.HeapObjects)
        fmt.Printf("\tHeapInuse = %v B", m.HeapInuse)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

