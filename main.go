package main

import (
	"awesomeProject/SimpleProxy"
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

// TODO: Maybe add a flags system to maybe set the config location and preferred logs location
func main() {
	fmt.Println("( WarpRelay Proxy )")
	wg.Add(4)
	go func() {
		defer wg.Done()
		SimpleProxy.SetupLogger()
	}()
	Spinner(2)
	// Proxy goroutine
	go func() {
		defer wg.Done()
		SimpleProxy.StartProxy()
	}()
	wg.Wait()
}
