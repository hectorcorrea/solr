package parser

// StringParser - parses the documents as map of maps
type StringParser struct {
}

// Parse - parses the pure document input from JSON
func (parser *StringParser) Parse(documents []byte) (interface{}, error) {
	var result interface{}
	result = (string)(documents)

	return result, nil
}
