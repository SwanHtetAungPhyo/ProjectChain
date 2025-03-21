package main

import (
	"fmt"
	"net"
	"sync"
)

func IsLive(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Port %d is LIVE (in use) ⚡\n", port)
	} else {
		_ = listener.Close()
		fmt.Printf("Port %d is FREE ✅\n", port)
	}

}
func main() {
	start := 8000
	end := 9000
	var wg sync.WaitGroup
	for i := start; i <= end; i++ {
		wg.Add(1)
		go IsLive(fmt.Sprintf(":%d", i), &wg)
	}
	wg.Wait()
	fmt.Printf("Port probing on the network done")
}
