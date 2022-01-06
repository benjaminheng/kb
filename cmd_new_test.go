package main

import "testing"

func TestTitleToFilename(t *testing.T) {
	var tests = []struct {
		name     string
		given    string
		expected string
	}{
		{
			"special characters replaced with dashes",
			"this is a title",
			"this-is-a-title.md",
		},
		{
			"runs of special characters replaced with single dash",
			"this is a    title",
			"this-is-a-title.md",
		},
		{
			"filename is lowercased",
			"Title",
			"title.md",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := titleToFilename(tt.given)
			if result != tt.expected {
				t.Errorf("(%+v): expected %+v, got %+v", tt.given, tt.expected, result)
			}

		})
	}
}
