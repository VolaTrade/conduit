package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	dataStreamer, shutdown, err := InitializeAndRun(ctx, "config.env")

	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(2)
	}

	defer shutdown()
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			cancel()
			time.Sleep(time.Duration(2 * time.Second))
			os.Exit(0)
		}
	}()

	dataStreamer.GetProcessCollectionState()
	dataStreamer.GenerateSocketListeningRoutines()

	go dataStreamer.ListenForDatabasePriveleges()
	go dataStreamer.RunSocketRoutines()
	go dataStreamer.ListenForExit(cancel)

	select {

	case <-c:
		println("Received OS SIGKILL")
		return

	case <-ctx.Done():
		println("Context shutdown")
		time.Sleep(time.Duration(2 * time.Second))
		return

	}

}
