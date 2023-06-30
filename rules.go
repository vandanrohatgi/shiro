// This package handles:
// - reading rules.yaml file
// - Associated functions
package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/itchyny/rassemble-go"
	"gopkg.in/yaml.v3"
)

// Rules & RuleConfig are used to read rules from rules.yaml
type RuleConfig struct {
	Path  string
	Rules map[string]Rules
}

type Rules struct {
	Body      string        `yaml:"body"`
	BodyRegex regexp.Regexp `yaml:"-"` // don't write these fields to rule file
	Headers   struct {
		Key        string        `yaml:"key"`
		KeyRegex   regexp.Regexp `yaml:"-"`
		Value      string        `yaml:"value"`
		ValueRegex regexp.Regexp `yaml:"-"`
	} `yaml:"headers"`
	Method      string        `yaml:"method"`
	MethodRegex regexp.Regexp `yaml:"-"`
	Description string        `yaml:"description"`
	Meta        string        `yaml:"meta"`
}

// Ruler contains methods associated with RuleConfig type
type Ruler interface {
	IngestRules()
	WriteRules()
}

// IngestRules reads the rules.yaml file and unmarshal them to ruleconfig instance
func (r *RuleConfig) IngestRules() {
	//rules := make(map[string]Rules)
	yamlFile, err := os.ReadFile(r.Path)
	if err != nil {
		log.Fatal("Error reading file", r.Path, err)
	}
	err = yaml.Unmarshal(yamlFile, &r.Rules)
	if err != nil {
		log.Fatal("Unable to extract rules", err)
	}
}

// GenerateRegex takes a list of strings and returns a regular expression string
func GenerateRegex(data []string) (string, error) {
	// TODO: create a custom library for better regex generation instead of current implemmentation
	pattern, err := rassemble.Join(data)
	if err != nil {
		return "", err
	}
	return pattern, nil
}

// WriteRules takes the RuleConfig struct and writes them to a yaml file.
// This function is used during monitoring mode.
func (r *RuleConfig) WriteRules() error {
	rules, err := yaml.Marshal(r.Rules)
	if err != nil {
		return err
	}
	err = os.WriteFile(r.Path, rules, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InspectBody(r *http.Request, rule Rules) Rules {
	// Generate Regex for body
	body, _ := io.ReadAll(r.Body)                // read request body
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore request body after reading it
	defer r.Body.Close()
	rule.Body, _ = GenerateRegex([]string{
		rule.Body,
		string(body),
	})
	return rule
}

func InspectHeaders(r *http.Request, rule Rules) Rules {
	// Generate Regex for headers
	for header, value := range r.Header {
		rule.Headers.Key, _ = GenerateRegex([]string{
			header,
			rule.Headers.Key,
		})
		rule.Headers.Value, _ = GenerateRegex([]string{
			strings.Join(value, ","),
			rule.Headers.Value,
		})
	}
	return rule
}

func InspectMethod(r *http.Request, rule Rules) Rules {
	//Generate Regex for Method
	rule.Method, _ = GenerateRegex([]string{
		r.Method,
		rule.Method,
	})
	return rule
}
