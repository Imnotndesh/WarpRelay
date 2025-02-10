package main

import (
	"awesomeProject/SimpleProxy"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

const helpText = `
USAGE:
	warp -c <path/to/custom/config> -> custom config file location
	warp -cv <path/to/custom/config> -> custom config file path and verbose output
	warp -h -> help text
	warp -v -> verbose log output`

var (
	wg          sync.WaitGroup
	arg         string
	defaultPath = ""
	customPath  string
)

const (
	ConfigDir      = "Config"
	ConfigLocation = "./Config/config.YAML"
)

func checkDefaultDir() error {
	// Checking if dir exists
	fmt.Println(" -> Checking default directories...")
	if _, err := os.Stat(ConfigDir); err != nil {
		if err = os.Mkdir(ConfigDir, os.ModePerm); err != nil {
			return errors.New("FAILED: Cannot generate Config directory")
		} else {
			fmt.Println(" -> Config directory created. PLease add a config file to the directory")
			os.Exit(0)
		}
	} else {
		// Check if config file in dir
		fmt.Println(" -> Checking for config file")
		if _, err = os.Stat(ConfigLocation); err != nil {
			return errors.New("cannot find config file: Please add a config file named 'config.YAML' to the Config directory")
		} else {
			return nil
		}
	}
	return nil
}
func checkCustomConfig(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("Cannot find custom config file", err)
		return
	}
}
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
	fmt.Println("\n ( WarpRelay Proxy )")
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("\n -> Using Default values")
		if err := checkDefaultDir(); err != nil {
			panic(err)
			return
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			SimpleProxy.SetupLogger()
		}()
		Spinner(2)
		go func() {
			defer wg.Done()
			SimpleProxy.StartProxy(defaultPath)
		}()
		wg.Wait()
	}
	arg = os.Args[1]
	switch arg {
	case "-v": // Verbose logs
		err := checkDefaultDir()
		if err != nil {
			panic(err)
			return
		}
		wg.Add(1)
		Spinner(2)
		go func() {
			defer wg.Done()
			SimpleProxy.StartProxy(defaultPath)
		}()
		wg.Wait()
	case "-cv": //Custom path and verbose logs
		if len(os.Args) < 3 || os.Args[2] == "" {
			fmt.Println("\n -> Please enter a valid path for the config")
			return
		}
		customPath = os.Args[2]
		checkCustomConfig(customPath)
		wg.Add(2)
		Spinner(2)
		go func() {
			defer wg.Done()
			SimpleProxy.StartProxy(customPath)
		}()
		wg.Wait()
	case "-c": // custom config default logging
		if len(os.Args) < 3 || os.Args[2] == "" {
			fmt.Println("\n -> Please enter a valid path for the config")
			return
		}
		customPath = os.Args[2]
		// Check if path is valid
		checkCustomConfig(customPath)
		wg.Add(2)
		go func() {
			defer wg.Done()
			SimpleProxy.SetupLogger()
		}()
		Spinner(2)
		go func() {
			defer wg.Done()
			SimpleProxy.StartProxy(customPath)
		}()
		wg.Wait()
	default:
		fallthrough
	case "-h": // Print help text
		fmt.Println(helpText)
	}
}
