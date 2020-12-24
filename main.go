package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	driver, shutdown, err := InitializeAndRun("config.env")

	if err != nil {
		fmt.Println(err)
		println("SHUTTTING DOWN")
		shutdown()
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer shutdown()
	c := make(chan os.Signal)

	var wg sync.WaitGroup

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()
	driver.RunDataStreamListenerRoutines(ctx, &wg)
	driver.Run(ctx, &wg, cancel)

	wg.Wait()

	println("LEAVING")
}
