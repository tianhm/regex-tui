package pcre

import (
	"go.elara.ws/pcre"
)

type PCRERegex struct {
	re *pcre.Regexp
}

func New(expr string) (*PCRERegex, error) {
	re, err := pcre.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &PCRERegex{re}, nil
}

func (regex *PCRERegex) FindAllStringIndex(s string, n int) [][]int {
	return regex.re.FindAllStringIndex(s, n)
}

func (regex *PCRERegex) FindStringIndex(s string) []int {
	return regex.re.FindStringIndex(s)
}
