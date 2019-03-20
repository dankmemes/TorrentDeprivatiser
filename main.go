package main

import (
	"fmt"
	"github.com/labstack/gommon/color"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var checkPre = color.Yellow("[") + color.Green("✓") + color.Yellow("]")
var tildPre = color.Yellow("[") + color.Green("~") + color.Yellow("]")
var crossPre = color.Yellow("[") + color.Red("✗") + color.Yellow("]")

var worker sync.WaitGroup
var count int

func main() {
	var err error

	// Parse arguments
	parseArgs(os.Args)

	// Check if input folder exist
	if _, err := os.Stat(arguments.Input); os.IsNotExist(err) {
		fmt.Println(crossPre +
			color.Yellow(" [") +
			color.Red(arguments.Input) +
			color.Yellow("] ") +
			color.Red("Invalid input folder."))
	}

	err = readTrackerList()
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(arguments.Input, WalkFile)
	if err != nil {
		log.Fatal(err)
	}

	worker.Wait()
}

func WalkFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	worker.Add(1)
	count++

	go func(path string) {
		err := work(path, &worker)

		if err != nil {
			fmt.Println(crossPre +
				color.Yellow(" [") +
				color.Red(path) +
				color.Yellow("] ") +
				color.Red("Error: ") +
				color.Yellow(err.Error()))
		}
	}(path)

	if count == arguments.Concurrency {
		worker.Wait()
		count = 0
	}

	return nil
}
