package re2

import (
	"regexp"
)

type RE2Regex struct {
	re *regexp.Regexp
}

func New(expr string) (*RE2Regex, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &RE2Regex{re}, nil
}

func (regex *RE2Regex) FindAllStringIndex(s string, n int) [][]int {
	return regex.re.FindAllStringIndex(s, n)
}

func (regex *RE2Regex) FindStringIndex(s string) []int {
	return regex.re.FindStringIndex(s)
}
