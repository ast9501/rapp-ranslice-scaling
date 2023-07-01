package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	DmaapIp             string `yaml:"DmaapIp"`
	DmaapPort           string `yaml:"DmaapPort"`
	DmaapTopic          string `yaml:"Topic"`
	CatalogueServiceUrl string `yaml:"CatalogueServiceUrl"`
}

func (c *Conf) ReadConf() *Conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
