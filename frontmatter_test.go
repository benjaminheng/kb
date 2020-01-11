package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseYAMLFrontMatter(t *testing.T) {
	t.Run("parses yaml from frontmatter", func(t *testing.T) {
		content := `---
title: this is a test
tags:
  - tag1
  - tag2
---

# Header

Content goes here
`
		expected := map[string]interface{}{
			"title": "this is a test",
			"tags":  []interface{}{"tag1", "tag2"},
		}
		result, err := ParseYAMLFrontMatter(strings.NewReader(content))
		if err != nil {
			t.Fatalf("expected nil error, got %+v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %+v, got %+v", expected, result)
		}
	})
	t.Run("returns ErrFrontMatterNotFound", func(t *testing.T) {
		content := `# Header

Content goes here
`
		_, err := ParseYAMLFrontMatter(strings.NewReader(content))
		if err != ErrFrontMatterNotFound {
			t.Fatalf("expected ErrFrontMatterNotFound, got %+v", err)
		}
	})
}
