//go:build ignore

// Command validate-schemas proves the shipped configuration schemas, examples,
// and the embedded Go contract agree. It fails when any example violates its
// schema enum/shape or when the Go strict decoder disagrees with the schema.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CGuiho/runx/pkg/config"
	"go.yaml.in/yaml/v3"
)

func bytesReader(data []byte) *bytes.Reader { return bytes.NewReader(data) }

func fatalf(format string, values ...any) {
	fmt.Fprintf(os.Stderr, "validate-schemas: "+format+"\n", values...)
	os.Exit(1)
}

type rawDocument struct {
	Schema string `yaml:"schema"`
	Agent  struct {
		Evolution struct {
			Upgrade string `yaml:"upgrade"`
			Issues  struct {
				Bugs         string `yaml:"bugs"`
				Improvements string `yaml:"improvements"`
				Reviews      string `yaml:"reviews"`
			} `yaml:"issues"`
		} `yaml:"evolution"`
	} `yaml:"agent"`
}

var allowedPolicies = map[string]bool{"disabled": true, "always-ask": true, "always-proceed": true, "": true}

func checkSchemaFile(path string) json.RawMessage {
	raw, err := os.ReadFile(path)
	if err != nil {
		fatalf("read %s: %v", path, err)
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		fatalf("parse %s: %v", path, err)
	}
	for _, key := range []string{"$schema", "$id", "title", "type", "additionalProperties", "properties"} {
		if _, ok := document[key]; !ok {
			fatalf("%s is missing required key %q", path, key)
		}
	}
	return raw
}

func checkExample(schemaPath, examplePath string, target any) {
	schemaRaw := checkSchemaFile(schemaPath)
	var schema struct {
		Properties struct {
			Agent struct {
				Ref string `json:"$ref"`
			} `json:"agent"`
		} `json:"properties"`
		Definitions map[string]any `json:"definitions"`
	}
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		fatalf("reparse %s: %v", schemaPath, err)
	}
	if len(schema.Definitions) == 0 || schema.Properties.Agent.Ref == "" {
		fatalf("%s lost its agent policy definitions", schemaPath)
	}

	exampleRaw, err := os.ReadFile(examplePath)
	if err != nil {
		fatalf("read %s: %v", examplePath, err)
	}
	var parsed rawDocument
	decoder := yaml.NewDecoder(bytesReader(exampleRaw))
	decoder.KnownFields(true)
	if err := decoder.Decode(&parsed); err != nil {
		fatalf("strict-decode %s: %v", examplePath, err)
	}
	for _, policy := range []string{
		parsed.Agent.Evolution.Upgrade,
		parsed.Agent.Evolution.Issues.Bugs,
		parsed.Agent.Evolution.Issues.Improvements,
		parsed.Agent.Evolution.Issues.Reviews,
	} {
		if !allowedPolicies[policy] {
			fatalf("%s uses policy value %q outside the schema enum", examplePath, policy)
		}
	}
	if parsed.Schema == "" {
		fatalf("%s is missing its version-pinned schema comment reference", examplePath)
	}

	// The Go runtime contract must accept the same bytes.
	if err := config.LoadDocument(exampleRaw, target); err != nil {
		fatalf("Go contract rejects %s: %v", examplePath, err)
	}
	fmt.Printf("[OK] %s agrees with %s\n", filepath.Base(examplePath), filepath.Base(schemaPath))
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatalf("resolve cwd: %v", err)
	}
	checkSchemaFile(filepath.Join(root, "schemas", "runx.schema.json"))
	checkSchemaFile(filepath.Join(root, "schemas", "runx.global.schema.json"))
	fmt.Println("[OK] both configuration schemas parse with complete policy contracts")

	checkExample(
		filepath.Join(root, "schemas", "runx.schema.json"),
		filepath.Join(root, "examples", "runx.yaml"),
		&config.ProjectConfig{},
	)
	checkExample(
		filepath.Join(root, "schemas", "runx.global.schema.json"),
		filepath.Join(root, "examples", "runx.global.yaml"),
		&config.GlobalConfig{},
	)
	fmt.Println("validate-schemas: all configuration contracts agree")
}
