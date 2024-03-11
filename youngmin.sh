name="youngmin"
role="normal"
port="4102"
peers="localhost:4101"
#peers="172.31.8.44:4101"
httpPort="9002"
key="f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7"
hostDataDir="/data/youngmin"
containerDataDir="/barreleye/barreldb/youngmin"

docker run -d -it --name ${name} --net host -v ${hostDataDir}:${containerDataDir} kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}