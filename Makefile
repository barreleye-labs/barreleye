build:
	go build -o ./bin/barreleye

run: build
	./bin/barreleye

genesis-node: build
	./bin/barreleye -nodeName=genesis-node

nayoung: build
	./bin/barreleye -nodeName=nayoung

youngmin: build
	./bin/barreleye -nodeName=youngmin

test:
	go test ./...