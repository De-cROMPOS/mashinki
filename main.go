package main

import (
	"fmt"
	"mashinki/logging"
	"mashinki/parser"
)

func main() {

	//Initializing logger
	logger, err := logging.NewFLogger("log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	url := "https://www.che168.com/dealer/166373/53093908.html"

	info, err := parser.GetCarInfo(url)
	if err != nil {
		logger.LogErrorF("Error while parsing car id: %v", err)
	}

	fmt.Println(info.String())
}
