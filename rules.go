package main

import (
	"log"
	"os"

	"github.com/itchyny/rassemble-go"
	"gopkg.in/yaml.v3"
)

type Rules struct {
	URI     string `yaml:"URI"`
	Body    string `yaml:"body"`
	Headers struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"headers"`
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

func GenerateRegex(data []string) (string, error) {
	pattern, err := rassemble.Join(data)
	if err != nil {
		return "", err
	}
	return pattern, nil
}
