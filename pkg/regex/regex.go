package regex

// Regex abstracts over compiled regular expressions for different engines.
type Regex interface {
	FindAllStringIndex(s string, n int) [][]int
	FindStringIndex(s string) []int
}
