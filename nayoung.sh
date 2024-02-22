name="nayoung"
role="n"
port="4101"
peer="172.30.1.5:4100"
httpPort="9001"
key="c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421"
docker run -it -p ${port}:${port} barreleye:latest /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peer=${peer} -http.port=${httpPort} -key=${key}