build:
	go build -o ./bin/barreleye

run: build
	./bin/barreleye

barreleye: build
	./bin/barreleye -name=barreleye -role=g -port=4100 -peer=localhost:4101 -http.port=9000 -key=a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9

nayoung: build
	./bin/barreleye -name=nayoung -role=n -port=4101 -peer=localhost:4100 -http.port=9001 -key=c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421

youngmin: build
	./bin/barreleye -name=youngmin -role=n -port=4102 -peer=localhost:4101 -http.port=9002 -key=f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7

test:
	go test ./...