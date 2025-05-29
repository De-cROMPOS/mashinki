package main

import (
	"fmt"
	"mashinki/Logger"
)

func main() {
	logger, err := logger.NewFLogger("log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.LogErrorF("Hello, I'm %s", "FLogger")
	fmt.Println("Hello World")
}
