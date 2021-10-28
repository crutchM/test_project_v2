package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		start := time.Now()
		PerformLinks()
		duration := time.Since(start)
		fmt.Printf("Elapsed time: %v", duration)
		time.Sleep(5000 * time.Second)
	}

}
