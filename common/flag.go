package common

import (
	"flag"
)

func ParseFlag() {
	flag.String("name", "", "node name")
	flag.String("role", "", "if genesis node -> g, if normal node -> n")
	flag.String("port", "", "port")
	flag.String("http.port", "", "http port")
	flag.String("peer", "", "peer")
	flag.String("key", "", "private key")
	flag.Parse()
}

func GetFlag(paramName string) string {
	return flag.Lookup(paramName).Value.(flag.Getter).Get().(string)
}
