package main

import (
	"awesomeProject/SimpleProxy"
	"awesomeProject/SimpleServer"
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	// Add yaml parsing for multiple endpoints
	wg.Add(2)
	go func() {
		defer SimpleServer.StopServer()
		defer wg.Done()
		fmt.Println("Starting simpleserver at :9080")
		SimpleServer.StartServer()
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Starting proxy server at :7080")
		SimpleProxy.StartProxy()
	}()
	wg.Wait()
}
