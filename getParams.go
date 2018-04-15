package solr

// GetParams represents the parameters used to get a single Solr document.
type GetParams struct {
	Q       string            // Typically "id:xyz"
	Fl      []string          // Fields to fetch from Solr.
	Options map[string]string // Options to pass straight to Solr (e.g. defType: "edismax")
}

// Creates a new GetParams object
func NewGetParams(q string, fl []string, options map[string]string) GetParams {
	params := GetParams{
		Q:       q,
		Fl:      fl,
		Options: options,
	}
	return params
}

func (params GetParams) toSolrQueryString() string {
	qs := ""
	qs += qsAdd("q", params.Q)
	qs += qsAddMany("fl", params.Fl)
	for k, v := range params.Options {
		qs += qsAdd(k, v)
	}
	return qs
}
