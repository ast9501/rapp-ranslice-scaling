package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	DmaapUrl            string `yaml:"DmaapUrl"`
	DmaapTopic          string `yaml:"DmaapTopic"`
	CatalogueServiceUrl string `yaml:"CatalogueServiceUrl"`
	NfvoUrl				string `yaml:"NfvoUrl"`
	VnfmUrl				string `yaml:"VnfmUrl"`
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
