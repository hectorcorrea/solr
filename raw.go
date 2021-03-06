package solr

import "encoding/json"

// The *Raw structs are used to unmarshall the JSON from Solr
// via Go's built-in functions. They are not exposed outside
// the solr package.
type headerRaw struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
	// Use interface{} because some params are strings and
	// others (e.g. fq) are arrays of strings.
	Params map[string]interface{} `json:"params"`
}

type dataRaw struct {
	NumFound  int           `json:"numFound"`
	Start     int           `json:"start"`
	Documents []documentRaw `json:"docs"`
}

// Just as it comes from Solr
type documentRaw map[string]interface{}

// fieldname: [value1, value2, ...]
type highlightRow map[string][]string

type errorRaw struct {
	Trace string `json:"trace"`
	Code  int    `json:"code"`
}

type responseRaw struct {
	Header       headerRaw               `json:"responseHeader"`
	Data         dataRaw                 `json:"response"`
	Error        errorRaw                `json:"error"`
	FacetCounts  facetCountsRaw          `json:"facet_counts"`
	Highlighting map[string]highlightRow `json:"highlighting"`
	Raw          string                  `json:"raw"`
}

func NewResponseRaw(rawBytes []byte) (responseRaw, error) {
	var response responseRaw
	err := json.Unmarshal([]byte(rawBytes), &response)
	if err != nil {
		return response, err
	}
	response.Raw = string(rawBytes)
	return response, nil
}

type facetCountsRaw struct {
	Queries interface{}              `json:"facet_queries"`
	Fields  map[string][]interface{} `json:"facet_fields"`
}
