package main

import (
	"awesomeProject/SimpleProxy"
	"awesomeProject/SimpleServer"
	"fmt"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func Spinner(seconds int) {
	frames := []string{".", "..", "...", "...."}
	i := 0
	endTime := time.Now().Add(time.Duration(seconds) * time.Second)

	for time.Now().Before(endTime) {
		fmt.Printf("\r %s ", frames[i])
		i = (i + 1) % len(frames)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("( WarpRelay Proxy )")
	wg.Add(4)
	go func() {
		defer wg.Done()
		SimpleProxy.SetupLogger()
	}()
	Spinner(2)
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
	// Proxy goroutine
	go func() {
		defer wg.Done()
		SimpleProxy.StartProxy()
	}()
	wg.Wait()
}
