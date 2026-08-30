package main

import (
	"reflect"
	"testing"
)

func TestTranslateArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "bare rx",
			input:    []string{},
			expected: []string{"list"},
		},
		{
			name:     "version flag short",
			input:    []string{"-v"},
			expected: []string{"-v"},
		},
		{
			name:     "version flag long",
			input:    []string{"--version"},
			expected: []string{"--version"},
		},
		{
			name:     "help flag short",
			input:    []string{"-h"},
			expected: []string{"-h"},
		},
		{
			name:     "help flag long",
			input:    []string{"--help"},
			expected: []string{"--help"},
		},
		{
			name:     "help-tree flag",
			input:    []string{"--help-tree"},
			expected: []string{"--help-tree"},
		},
		{
			name:     "help-tree-depth flag with value",
			input:    []string{"--help-tree-depth", "2"},
			expected: []string{"--help-tree-depth", "2"},
		},
		{
			name:     "help-tree-depth flag with equals",
			input:    []string{"--help-tree-depth=2"},
			expected: []string{"--help-tree-depth=2"},
		},
		{
			name:     "help-tree-global-flags",
			input:    []string{"--help-tree-global-flags"},
			expected: []string{"--help-tree-global-flags"},
		},
		{
			name:     "help-docs flag",
			input:    []string{"--help-docs"},
			expected: []string{"--help-docs"},
		},
		{
			name:     "color flag single",
			input:    []string{"--color"},
			expected: []string{"--color"},
		},
		{
			name:     "list with format json",
			input:    []string{"--format", "json"},
			expected: []string{"list", "--format", "json"},
		},
		{
			name:     "list with cwd",
			input:    []string{"--cwd", "/tmp"},
			expected: []string{"list", "--cwd", "/tmp"},
		},
		{
			name:     "list with config",
			input:    []string{"--config", "custom.yaml"},
			expected: []string{"list", "--config", "custom.yaml"},
		},
		{
			name:     "list with verbose",
			input:    []string{"--verbose"},
			expected: []string{"list", "--verbose"},
		},
		{
			name:     "run with selector",
			input:    []string{"build"},
			expected: []string{"run", "build"},
		},
		{
			name:     "run with numeric index",
			input:    []string{"1"},
			expected: []string{"run", "1"},
		},
		{
			name:     "run with selector and child args",
			input:    []string{"build", "--foo", "bar"},
			expected: []string{"run", "build", "--foo", "bar"},
		},
		{
			name:     "run with selector and delimiter child args",
			input:    []string{"build", "--", "--dry-run"},
			expected: []string{"run", "build", "--", "--dry-run"},
		},
		{
			name:     "run with options before selector",
			input:    []string{"--cwd", "/tmp", "deploy"},
			expected: []string{"run", "--cwd", "/tmp", "deploy"},
		},
		{
			name:     "run with dry-run before selector",
			input:    []string{"--dry-run", "1"},
			expected: []string{"run", "--dry-run", "1"},
		},
		{
			name:     "run with yes before selector",
			input:    []string{"--yes", "deploy"},
			expected: []string{"run", "--yes", "deploy"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := translateArgs(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("translateArgs(%v) = %v, want %v", tc.input, actual, tc.expected)
			}
		})
	}
}
