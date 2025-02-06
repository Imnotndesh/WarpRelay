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
	wg.Add(2)
	go func() {
		defer SimpleServer.StopServer()
		defer wg.Done()
		SimpleServer.StartServer()
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
