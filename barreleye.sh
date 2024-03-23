name="barreleye"
role="genesis"
port="4100"
peers="none"
httpPort="9000"
key=
hostDataDir="/data/barreleye"
containerDataDir="/barreleye/barreldb/barreleye"

docker run -d -it --name ${name} --net host -v ${hostDataDir}:${containerDataDir} kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peers=${peers} -http.port=${httpPort} -key=${key}