package main

import (
	"fmt"

	config "github.com/ast9501/rapp-ranslice-scaling/internal"
)

func main() {
	var c config.Conf
	c.ReadConf()

	fmt.Println(c.DmaapIp)
}
