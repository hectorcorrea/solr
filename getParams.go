package solr

type GetParams struct {
	Q       string
	Fl      []string
	Options map[string]string
}

// NewSearchParams from a query string
// 	`q` is typically "id:xyz"
// 	`fl` list of fields to fetch
// 	`options` to pass to Solr (e.g. defType: "edismax", wt: "json")
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
	qs += QsAdd("q", params.Q)
	qs += QsAddMany("fl", params.Fl)
	for k, v := range params.Options {
		qs += QsAdd(k, v)
	}
	return qs
}
