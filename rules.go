package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Rules struct {
	URI         string `yaml:"URI"`
	Body        string `yaml:"body"`
	Headers     string `yaml:"headers"`
	Method      string `yaml:"method"`
	Description string `yaml:"description"`
	Meta        string `yaml:"meta"`
}

type RuleConfig struct {
	Path       string
	RulesArray []Rules
}

type RuleMethods interface {
	IngestRules()
	PrintRules()
	GetInstance() RuleConfig
}

func (r *RuleConfig) IngestRules() {
	yamlFile, err := os.ReadFile(r.Path)
	if err != nil {
		log.Fatal("Error reading file", r.Path, err)
	}
	err = yaml.Unmarshal(yamlFile, &r.RulesArray)
	if err != nil {
		log.Fatal("Unable to extract rules", err)
	}
}

func (r *RuleConfig) PrintRules() {
	log.Println(r.RulesArray)
}
