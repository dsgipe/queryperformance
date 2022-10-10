package main

import (
	"fmt"
	"os"
	"timescale/queryperformance/worker"
)

func main() {
	var (
		input [][]string
		err   error
	)
	if len(os.Args) == 3 && os.Args[1] == "-f" {
		input, err = worker.ImportCsv(os.Args[2])
	} else {
		input, err = worker.ProcessInput(os.Stdin)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = worker.EvaluateTimescaleQuerySpeeds(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
