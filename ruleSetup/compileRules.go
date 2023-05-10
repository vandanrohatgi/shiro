package rules

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Rules []struct {
	Path         string `yaml:"path"`
	AllowPattern string `yaml:"allowPattern"`
	Method       string `yaml:"method"`
	Description  string `yaml:"description"`
	Meta         string `yaml:"meta"`
}

func PrintRules() {
	yamlFile, err := os.ReadFile("rules.yaml")
	if err != nil {
		log.Fatal("Error reading rules.yaml file", err)
	}
	var rule Rules
	err = yaml.Unmarshal(yamlFile, &rule)
	if err != nil {
		log.Fatal("Unable to extract rules", err)
	}
	log.Println(rule)
}
