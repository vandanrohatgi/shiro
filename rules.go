package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Rules struct {
	Path                string `yaml:"path"`
	AllowPatternURI     string `yaml:"allowPatternURI"`
	AllowPatternBody    string `yaml:"allowPatternBody"`
	AllowPatternHeaders string `yaml:"allowPatternHeaders"`
	Method              string `yaml:"method"`
	Description         string `yaml:"description"`
	Meta                string `yaml:"meta"`
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
