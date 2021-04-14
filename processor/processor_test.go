package processor_test

import (
	p "github.com/ashkarin/httpresphasher/processor"
	"testing"
)

func appendHello(value p.Value) p.Result {
	return p.Result{
		Input:  value,
		Output: p.Value(string(value) + "Hello"),
		Error:  nil,
	}
}

func generateValues() (map[string]string, chan p.Value) {
	data := []string{"cat", "dog", "rat", "truck", ""}
	reference := make(map[string]string)
	values := make(chan p.Value, len(data))

	for _, v := range data {
		values <- p.Value(v)
		reference[v] = v + "Hello"
	}
	close(values)

	return reference, values
}

func TestProcess(t *testing.T) {
	var nWorkers uint
	for nWorkers = 1; nWorkers < 10; nWorkers++ {
		reference, valuesQueue := generateValues()
		results, _ := p.Process(valuesQueue, appendHello, nWorkers, 1)
		for result := range results {
			if reference[string(result.Input)] != string(result.Output) {
				t.Errorf("Processor returned wrong result: %v != %v", result.Input, reference[string(result.Input)])
			}
		}
	}

	nWorkers = 0
	_, valuesQueue := generateValues()
	results, err := p.Process(valuesQueue, appendHello, nWorkers, 1)
	if err == nil {
		t.Error("Processor should return an error when nWorkers=0")
	}
	if results != nil {
		t.Error("Processor should not return a generator when nWorkers=0")
	}
}
