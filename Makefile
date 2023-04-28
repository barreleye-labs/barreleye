build:
		go build  -o ./bin/barreleye

run: build
		./bin/barreleye
		
test: 
		go test ./...