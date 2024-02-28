name="youngmin"
role="normal"
port="4102"
peers="172.31.8.44:4101"
httpPort="9002"
key="f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7"
docker run -d -it --name ${name} -p ${port}:${port} -p ${httpPort}:${httpPort} -v /data:/data kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}