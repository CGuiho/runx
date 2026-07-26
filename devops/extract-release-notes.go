//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	version := flag.String("version", "", "exact semantic version heading to extract")
	input := flag.String("input", "CHANGELOG.md", "changelog path")
	output := flag.String("output", "release-notes.md", "release notes path")
	flag.Parse()
	if *version == "" {
		fatalf("--version is required")
	}
	data, err := os.ReadFile(*input)
	if err != nil {
		fatalf("read changelog: %v", err)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	heading := regexp.MustCompile(`(?m)^## ` + regexp.QuoteMeta(*version) + ` - [^\n]+\n`)
	location := heading.FindStringIndex(text)
	if location == nil {
		fatalf("release notes section for %s was not found", *version)
	}
	end := len(text)
	if next := regexp.MustCompile(`(?m)^## `).FindStringIndex(text[location[1]:]); next != nil {
		end = location[1] + next[0]
	}
	section := strings.TrimSpace(text[location[0]:end]) + "\n"
	if strings.TrimSpace(section) == strings.TrimSpace(text[location[0]:location[1]]) {
		fatalf("release notes section for %s is empty", *version)
	}
	if err := os.WriteFile(*output, []byte(section), 0o644); err != nil {
		fatalf("write release notes: %v", err)
	}
}
func fatalf(format string, values ...any) { fmt.Fprintf(os.Stderr, format+"\n", values...); os.Exit(1) }
