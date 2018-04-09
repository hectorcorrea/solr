package solr

import (
	"net/url"
)

const defaultRows = 10

type SearchParams struct {
	Q             string
	Fl            []string
	Rows          int
	Start         int
	FilterQueries FilterQueries
	Facets        Facets
	Options       map[string]string
}

// NewSearchParams from a query string. Will pick up several known
// values from the query string (e.g. q, rows, start, fq)
//
// 	`qs` is typically req.URL.Query()
// 	`options` to pass to Solr (e.g. defType: "edismax")
// 	`facets` to request from Solr (e.g. fieldName: "Field Name")
func NewSearchParamsFromQs(qs url.Values, options map[string]string,
	facets map[string]string) SearchParams {

	params := SearchParams{
		Q:             QsGet(qs, "q", "*"),
		Rows:          QsGetInt(qs, "rows", defaultRows),
		Start:         QsGetInt(qs, "start", 0),
		FilterQueries: NewFilterQueries(qs["fq"]),
		Options:       options,
		Facets:        NewFacets(facets),
	}
	return params
}

// NewSearchParams from a search string. You cannot set filter queries or
// other parameters with this option. But you can set them on the returned
// object.
//
// 	`q` value to pass to Solr's q parameter.
// 	`options` to pass to Solr (e.g. defType: "edismax")
// 	`facets` to request from Solr (e.g. fieldName: "Field Name")
func NewSearchParams(q string, options map[string]string,
	facets map[string]string) SearchParams {

	params := SearchParams{
		Q:       q,
		Options: options,
		Facets:  NewFacets(facets),
	}
	params.Rows = defaultRows
	return params
}

func (params SearchParams) toSolrQueryString() string {
	qs := ""
	qs += QsAddDefault("q", params.Q, "*")
	qs += QsAddMany("fl", params.Fl)
	qs += params.FilterQueries.toQueryString()
	qs += params.Facets.toQueryString()

	if params.Start > 0 {
		qs += QsAddInt("start", params.Start)
	}

	if params.Rows != defaultRows {
		qs += QsAddInt("rows", params.Rows)
	}

	for k, v := range params.Options {
		qs += QsAdd(k, v)
	}
	return qs
}
