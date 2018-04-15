package solr

type GetParams struct {
	Q       string
	Fl      []string
	Options map[string]string
}

// NewGetParams
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
	qs += qsAdd("q", params.Q)
	qs += qsAddMany("fl", params.Fl)
	for k, v := range params.Options {
		qs += qsAdd(k, v)
	}
	return qs
}
