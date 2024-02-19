package common

import (
	"flag"
)

func GetFlag() (nodeName string) {
	nodeName = ""
	flag.StringVar(&nodeName, "nodeName", "", "node name")
	flag.Parse()
	return
}
