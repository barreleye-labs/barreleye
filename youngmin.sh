name="youngmin"
role="normal"
port="4102"
peers="172.30.1.5:4101"
httpPort="9002"
key="f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7"
docker run -it -p ${port}:${port} -p ${httpPort}:${httpPort} --name ${name} kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}