build:
	go build -o ./bin/barreleye

run: build
	./bin/barreleye

barreleye: build
	./bin/barreleye -name=barreleye -role=genesis -port=4100 -peers=none -http.port=9000 -key=

nayoung: build
	./bin/barreleye -name=nayoung -role=normal -port=4101 -peers=localhost:4100 -http.port=9001 -key=

youngmin: build
	./bin/barreleye -name=youngmin -role=normal -port=4102 -peers=localhost:4101 -http.port=9002 -key=

test:
	go test ./...