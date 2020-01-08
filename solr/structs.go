package solr

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

type dataRawInterface struct {
	NumFound  int                    `json:"numFound"`
	Start     int                    `json:"start"`
	Documents []documentRawInterface `json:"docs"`
}

type documentRawInterface interface{}

// Just as it comes from Solr
type documentRaw map[string]interface{}

// fieldname: [value1, value2, ...]
type highlightRow map[string][]string

type errorRaw struct {
	Trace string `json:"trace"`
	Code  int    `json:"code"`
}

type responseRaw struct {
	Header      headerRaw      `json:"responseHeader"`
	Data        dataRaw        `json:"response"`
	Error       errorRaw       `json:"error"`
	FacetCounts facetCountsRaw `json:"facet_counts"`
}

type facetCountsRaw struct {
	Queries interface{}              `json:"facet_queries"`
	Fields  map[string][]interface{} `json:"facet_fields"`
}

//Response - response data from solr instance
type Response struct {
	Status   int64        `json:"status,omitempty"`
	QTime    int64        `json:"Qtime,omitempty"`
	NumFound int64        `json:"numFound,omitempty"`
	MaxScore float64      `json:"maxScore,omitempty"`
	Docs     interface{}  `json:"Docs,omitempty"`
	Facets   []facetField `json:"Facets,omitempty"`
}

type facetField struct {
	name string
	list []facetValue
}
type facetValue struct {
	name  string
	value int64
}
