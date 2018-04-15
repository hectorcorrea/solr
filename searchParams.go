package solr

import (
	"net/url"
)

const defaultRows = 10

// SearchParams represents the parameters used to issue a
// search in Solr.
type SearchParams struct {
	Q             string
	Fl            []string
	Rows          int
	Start         int
	FilterQueries filterQueries     // Values that will be passed as the fq parameter.
	Facets        Facets            // Facets to request from Solr.
	Options       map[string]string // Options to pass straight to Solr (e.g. defType: "edismax")
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
		FilterQueries: newFilterQueries(qs["fq"]),
		Options:       options,
		Facets:        newFacets(facets),
	}
	return params
}

// NewSearchParams from a search string.
func NewSearchParams(q string, options map[string]string,
	facets map[string]string) SearchParams {

	params := SearchParams{
		Q:       q,
		Options: options,
		Facets:  newFacets(facets),
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
