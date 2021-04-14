package processor

import (
	"fmt"
	"sync"
)

type Value string

var NoValue Value = ""

type Result struct {
	Input  Value
	Output Value
	Error  error
}

type Processor func(Value) Result

// Process processes values with processor using nWorkers and returns
// an unordered generator of Result
func Process(values chan Value, processor Processor, nWorkers uint, queueSize uint) (<-chan Result, error) {
	if nWorkers == 0 {
		return nil, fmt.Errorf("nWorkers should be greater or equal to 1")
	}

	results := make(chan Result)
	workersQueues := make(chan (chan<- Value), nWorkers)
	workersResults := make(chan (<-chan Result), nWorkers)
	wg := new(sync.WaitGroup)
	wg.Add(int(nWorkers))

	// spawn workers
	go func() {
		var i uint
		for i = 0; i < nWorkers; i++ {
			workerQueue := make(chan Value, queueSize)
			workersQueues <- workerQueue
			workersResults <- worker(workerQueue, processor)
		}
	}()

	// aggregate results
	go func() {
		for {
			workerResults, opened := <-workersResults
			if !opened {
				return
			}

			result, opened := <-workerResults
			if !opened {
				wg.Done()
				continue
			}

			workersResults <- workerResults
			results <- result
		}
	}()

	// finalize
	go func() {
		wg.Wait()
		close(workersQueues)
		close(workersResults)
		close(results)
	}()

	// distribute values
	go func() {
		for value := range values {
			workerQueue := <-workersQueues
			workerQueue <- value
			workersQueues <- workerQueue
		}

		for workerQueue := range workersQueues {
			close(workerQueue)
		}
	}()

	return results, nil
}

func worker(values <-chan Value, processor Processor) <-chan Result {
	results := make(chan Result)
	go func() {
		for {
			value, opened := <-values
			if !opened {
				close(results)
				return
			}
			results <- processor(value)
		}
	}()
	return results
}
