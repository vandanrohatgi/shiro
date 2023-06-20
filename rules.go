// This package handles:
// - reading rules.yaml file
// - Assocaited functions
package main

import (
	"log"
	"os"

	"github.com/itchyny/rassemble-go"
	"gopkg.in/yaml.v3"
)

// Rules & RuleConfig are used to read rules from rules.yaml
type RuleConfig struct {
	Path       string
	RulesArray []Rules
}

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

// Ruler contains methods associated with RuleConfig type
type Ruler interface {
	IngestRules()
	PrintRules()
}

// IngestRules reads the rules.yaml file and unmarshal them to ruleconfig instance
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

// PrintRules is used for debugging purposes
// Prints all the rules to stdout which were ingested on initializaition
func (r *RuleConfig) PrintRules() {
	log.Println(r.RulesArray)
}

// GenerateRegex takes a list of strings and returns a regular expression string
func GenerateRegex(data []string) (string, error) {
	pattern, err := rassemble.Join(data)
	if err != nil {
		return "", err
	}
	return pattern, nil
}
