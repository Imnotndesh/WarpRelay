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

func main() {
	// Launch a goroutine for the reverse proxy and the simple server
	wg.Add(3)
	go func() {
		defer SimpleServer.StopServer()
		defer wg.Done()
		SimpleServer.StartServer(":9080", "Hello World at 9080")
	}()
	go func() {
		defer wg.Done()
		defer SimpleServer.StopServer()
		SimpleServer.StartServer(":9081", "Hello World at 9081")
	}()
	time.Sleep(2 * time.Second)
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
