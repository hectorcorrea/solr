package solr

import (
	"net/url"
)

const defaultRows = 10

// SearchParams represents the parameters used to issue a
// search in Solr. Q, Fl, Rows, and Start map to the Solr
// equivalent.
//
// FilterQueries represents the values that will be passed
// to Solr as the fq parameter.
//
// Facets is an array with the facets to request from Solr.
//
// Options is map with the options to pass straight to Solr
// (e.g. defType: "edismax")
type SearchParams struct {
	Q             string
	Fl            []string
	Rows          int
	Start         int
	FilterQueries FilterQueries
	Facets        Facets
	Options       map[string]string
}

// NewSearchParams from a query string. This method will automatically
// pickup several known parameters from the query string (q, rows,
// start, and fq).
//
// qs typically an instance of req.URL.Query() from a web handler.
func NewSearchParamsFromQs(qs url.Values, options map[string]string,
	facets map[string]string) SearchParams {

	params := SearchParams{
		Q:             qsGet(qs, "q", "*"),
		Rows:          qsGetInt(qs, "rows", defaultRows),
		Start:         qsGetInt(qs, "start", 0),
		FilterQueries: NewFilterQueries(qs["fq"]),
		Options:       options,
		Facets:        NewFacets(facets),
	}
	return params
}

// NewSearchParams from a search string.
func newSearchParams(q string, options map[string]string,
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
	qs += qsAddDefault("q", params.Q, "*")
	qs += qsAddMany("fl", params.Fl)
	qs += params.FilterQueries.toQueryString()
	qs += params.Facets.toQueryString()

	if params.Start > 0 {
		qs += qsAddInt("start", params.Start)
	}

	if params.Rows != defaultRows {
		qs += qsAddInt("rows", params.Rows)
	}

	for k, v := range params.Options {
		qs += qsAdd(k, v)
	}
	return qs
}
