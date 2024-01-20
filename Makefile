build:
	go build -o ./bin/barreleye

run: build
	./bin/barreleye

node1: build
	./bin/barreleye -nodeName=node1

node2: build
	./bin/barreleye -nodeName=node2

test:
	go test ./...