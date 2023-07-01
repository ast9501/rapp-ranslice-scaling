package main

import (
	"fmt"

	internal "github.com/ast9501/rapp-ranslice-scaling/internal"
)

func main() {
	var c internal.Conf
	c.ReadConf()
	internal.Register(c.CatalogueServiceUrl)

	fmt.Println(c.DmaapIp)
}
