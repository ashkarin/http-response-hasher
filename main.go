package main

import (
	"flag"
	"fmt"
	"http-response-hasher/app"
	"os"
)

func main() {
	nWorkers := flag.Uint("parallel", 10, "number of workers")
	flag.Parse()

	urls := flag.Args()
	if len(urls) < 1 {
		fmt.Printf("No URLs specified")
		os.Exit(1)
	}

	results, err := app.ProcessUrls(urls, *nWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for result := range results {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "%s %v\n", result.Input, result.Error)
		} else {
			fmt.Printf("%s %s\n", result.Input, result.Output)
		}
	}
}
