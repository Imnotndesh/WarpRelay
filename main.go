package main

import (
	"awesomeProject/SimpleProxy"
	"awesomeProject/SimpleServer"
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
		SimpleServer.StartServer()
	}()
	go func() {
		defer wg.Done()
		SimpleProxy.StartProxy()
	}()
	wg.Wait()
}
