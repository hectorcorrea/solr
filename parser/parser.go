package parser

// DocumentParser - a generic document parser
type DocumentParser interface {

	// Parse - parses the pure document input from JSON
	Parse(documents []byte) (interface{}, error)
}
