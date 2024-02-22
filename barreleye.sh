name="barreleye"
role="g"
port="4100"
peer="172.30.1.5:4101"
httpPort="9000"
key="a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9"
docker run -it -p ${port}:${port} barreleye:latest /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peer=${peer} -http.port=${httpPort} -key=${key}