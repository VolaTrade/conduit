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

	dataStreamer, shutdown, err := InitializeAndRun("config.env")

	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer shutdown()
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {

		select {

		case <-c:
			cancel()
			time.Sleep(time.Duration(2 * time.Second))
			os.Exit(1)
		}
	}()

	if err := dataStreamer.InsertPairsFromBinanceToCache(); err != nil {
		
		panic(err)
	}
	dataStreamer.GenerateSocketListeningRoutines(ctx)

	go dataStreamer.ListenForDatabasePriveleges(ctx)
	go dataStreamer.RunSocketRoutines(ctx)
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
