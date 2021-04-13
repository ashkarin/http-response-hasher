init:
	go mod tidy

build: init
	go build

test: init
	ginkgo -r
