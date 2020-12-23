package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	driver, end, err := InitializeAndRun("config.env")

	if err != nil {
		fmt.Println(err)
		end()
		os.Exit(2)
	}

	defer end()
	c := make(chan os.Signal)
	quit := make(chan bool)
	var wg sync.WaitGroup
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quit <- true
		os.Exit(1)
	}()
	driver.RunDataStreamListenerRoutines(&wg, quit)
	driver.Run(&wg)

	wg.Wait()
}
