build:
	go build -o ./bin/barreleye

run: build
	./bin/barreleye

genesis-node: build
	./bin/barreleye -nodeName=genesis-node

wayne: build
	./bin/barreleye -nodeName=wayne

usi: build
	./bin/barreleye -nodeName=usi

test:
	go test ./...