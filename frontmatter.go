package main

import (
	"bufio"
	"errors"
	"io"

	"gopkg.in/yaml.v2"
)

var ErrFrontMatterNotFound = errors.New("front matter not found")

func ParseYAMLFrontMatter(input io.Reader) (result map[string]interface{}, err error) {
	s := bufio.NewScanner(input)
	var line int64
	var isInFrontMatter bool
	var frontMatter string
	for s.Scan() {
		line++
		lineContent := s.Text()
		if line == 1 && lineContent == "---" {
			isInFrontMatter = true
			continue
		}
		if isInFrontMatter && lineContent == "---" {
			break
		}
		if isInFrontMatter {
			frontMatter += lineContent + "\n"
			continue
		}
	}
	if len(frontMatter) > 0 {
		err = yaml.Unmarshal([]byte(frontMatter), &result)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrFrontMatterNotFound
	}
	return result, nil
}
