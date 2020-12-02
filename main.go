package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	driver, err := InitializeAndRun("config.env")

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	c := make(chan os.Signal)
	quit := make(chan bool)

	var wg sync.WaitGroup
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		quit <- true
		os.Exit(1)
	}()
	driver.InitService()
	driver.RunListenerRoutines(&wg, quit)
	driver.Run(&wg)

	wg.Wait()
}
