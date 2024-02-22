name="youngmin"
role="n"
port="4102"
peer="172.30.1.5:4101"
httpPort="9002"
key="f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7"
docker run -it -p ${port}:${port} barreleye:latest /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peer=${peer} -http.port=${httpPort} -key=${key}