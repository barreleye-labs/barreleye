name="nayoung"
role="normal"
port="4101"
peers="localhost:4100"
httpPort="9001"
key=
hostDataDir="/data/nayoung"
containerDataDir="/barreleye/barreldb/nayoung"

docker run -d -it --name ${name} --net host -v ${hostDataDir}:${containerDataDir} kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}