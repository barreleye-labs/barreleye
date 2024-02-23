name="nayoung"
role="normal"
port="4101"
peers="172.30.1.5:4100"
httpPort="9001"
key="c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421"
docker run -it -p ${port}:${port} -p ${httpPort}:${httpPort} --name ${name} kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}