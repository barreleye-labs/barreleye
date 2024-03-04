name="nayoung"
role="normal"
port="4101"
peers="localhost:4100"
#peers="172.31.8.44:4100"
httpPort="9001"
key="c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421"
docker run -d -it --name ${name} --net host -v /data:/data kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}
#docker run -d -it --name ${name} -p ${port}:${port} -p ${httpPort}:${httpPort} -v /data:/data kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}