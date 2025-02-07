package main

import (
	"awesomeProject/SimpleProxy"
	"awesomeProject/SimpleServer"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

// TODO: Add a custom logger that writes to file in /logs dir
func main() {
	wg.Add(3)
	// Simple servers to test out proxy
	go func() {
		defer SimpleServer.StopServer()
		defer wg.Done()
		SimpleServer.StartServer(":9080", "Hello World at 9080", false)
	}()
	go func() {
		defer wg.Done()
		defer SimpleServer.StopServer()
		SimpleServer.StartServer(":9081", "Hello World at 9081", true)
	}()
	time.Sleep(2 * time.Second)

	// Proxy goroutine
	go func() {
		defer wg.Done()
		SimpleProxy.StartProxy()
	}()
	wg.Wait()

	// For debug purposes
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	SimpleServer.StartServer()
	//}()
	//SimpleProxy.StartProxy()
}
