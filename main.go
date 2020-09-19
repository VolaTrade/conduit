package main

import (
	"fmt"
	"os"
)

func main() {

	driver, err := InitializeAndRun("config.env")

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	driver.Run()
}
