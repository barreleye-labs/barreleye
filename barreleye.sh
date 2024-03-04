name="barreleye"
role="genesis"
port="4100"
peers="none"
httpPort="9000"
key="a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9"
docker run -d -it --name ${name} --net host -v /data:/data kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}
#docker run -d -it --name ${name} -p ${port}:${port} -p ${httpPort}:${httpPort} -v /data:/data kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}