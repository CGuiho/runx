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
			name:     "bare rxs runs payload smart",
			input:    []string{},
			expected: []string{"smart"},
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
			name:     "list options forward through smart",
			input:    []string{"--format", "json"},
			expected: []string{"smart", "--format", "json"},
		},
		{
			name:     "list with cwd",
			input:    []string{"--cwd", "/tmp"},
			expected: []string{"smart", "--cwd", "/tmp"},
		},
		{
			name:     "list with config",
			input:    []string{"--config", "custom.yaml"},
			expected: []string{"smart", "--config", "custom.yaml"},
		},
		{
			name:     "list with verbose",
			input:    []string{"--verbose"},
			expected: []string{"smart", "--verbose"},
		},
		{
			name:     "selector forwards through smart",
			input:    []string{"build"},
			expected: []string{"smart", "build"},
		},
		{
			name:     "numeric index forwards through smart",
			input:    []string{"1"},
			expected: []string{"smart", "1"},
		},
		{
			name:     "selector with child args",
			input:    []string{"build", "--foo", "bar"},
			expected: []string{"smart", "build", "--foo", "bar"},
		},
		{
			name:     "selector with delimiter child args",
			input:    []string{"build", "--", "--dry-run"},
			expected: []string{"smart", "build", "--", "--dry-run"},
		},
		{
			name:     "options before selector",
			input:    []string{"--cwd", "/tmp", "deploy"},
			expected: []string{"smart", "--cwd", "/tmp", "deploy"},
		},
		{
			name:     "dry-run before selector",
			input:    []string{"--dry-run", "1"},
			expected: []string{"smart", "--dry-run", "1"},
		},
		{
			name:     "yes before selector",
			input:    []string{"--yes", "deploy"},
			expected: []string{"smart", "--yes", "deploy"},
		},
		{
			name:     "list or run translation stays in the payload",
			input:    []string{"--format", "json", "build"},
			expected: []string{"smart", "--format", "json", "build"},
		},
		{
			name:     "explicit dry-run false forwards through smart",
			input:    []string{"--dry-run=false", "deploy"},
			expected: []string{"smart", "--dry-run=false", "deploy"},
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
