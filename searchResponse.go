package solr

import "net/url"

// SearchResponse stores the results of a search, including
// the parameters used in the search, the resulting documents,
// and facets returned by the server. It also provides the URLs
// that could be submitted to NewSearchParamsFromQs() to re-execute
// the search without the Q parameter or to get the previous/next
// batch of results.
type SearchResponse struct {
	Params      SearchParams
	Q           string
	NumFound    int
	Start       int
	Rows        int
	Documents   []Document // Documents returned by Solr (including highlight information)
	Facets      Facets     // Facet information (field, title, and values)
	Url         string     // URL to execute this search
	UrlNoQ      string     // URL to execute this search without the Q parameter
	NextPageUrl string     // URL to get the next batch of results
	PrevPageUrl string     // URL to get the previous batch of results
	Raw         string
}

func newSearchResponse(params SearchParams, raw responseRaw) SearchResponse {
	r := SearchResponse{
		Params:    params,
		Q:         params.Q,
		NumFound:  raw.Data.NumFound,
		Start:     raw.Data.Start,
		Rows:      params.Rows,
		Documents: newDocumentFromSolrResponse(raw),
		Raw:       raw.Raw,
	}

	// Make sure the facets in the results are ordered according
	// to the facets in the request.
	unorderedFacets := r.facetsFromResponse(raw.FacetCounts)
	orderedFacets := Facets{}
	for _, facetDef := range params.Facets {
		facet, found := unorderedFacets.ForField(facetDef.Field)
		if found {
			orderedFacets = append(orderedFacets, facet)
		}
	}

	r.Facets = orderedFacets
	r.Url = r.toQueryString(r.Q, r.Start)
	r.UrlNoQ = r.toQueryString("", r.Start)
	r.NextPageUrl = r.toQueryString(r.Q, r.Start+r.Rows)
	r.PrevPageUrl = r.toQueryString(r.Q, r.Start-r.Rows)

	return r
}

func (r SearchResponse) toQueryString(q string, start int) string {
	qs := ""

	if q != "" {
		qs += qsAddRaw("q", q)
	}

	for _, facet := range r.Facets {
		for _, value := range facet.Values {
			if value.Active {
				valueEnc := url.QueryEscape(value.Text)
				qs += qsAddRaw("fq", facet.Field+"|"+valueEnc)
			}
		}
	}

	if start > 0 {
		qs += qsAddInt("start", start)
	}

	if r.Rows != 10 {
		qs += qsAddInt("rows", r.Rows)
	}
	return qs
}

// Creates a new Facets object from the raw FacetCounts from Solr.
// 	`counts` contains the facet data as reported by Solr.
func (r SearchResponse) facetsFromResponse(counts facetCountsRaw) Facets {
	facets := Facets{}
	for fieldName, tokens := range counts.Fields {

		facet := r.newFacet(fieldName)

		if len(tokens) > 0 {
			// Tokens is an array in the form [value1, count1, value2, count2].
			// Here we break it into an array of FacetValues that has specific
			// value and count properties.
			for i := 0; i < len(tokens); i += 2 {
				text := tokens[i].(string)
				count := int(tokens[i+1].(float64))
				// Mark the facet for this value as active if it is
				// present on the FilterQueries
				active := r.Params.FilterQueries.HasFieldValue(fieldName, text)
				facet.addValue(text, count, active)
			}
		} else {
			// If no data was returned for this field from Solr,
			// make sure to add it with a count of zero (if it was
			// on the filter queries.)
			for _, fqValue := range r.Params.FilterQueries.FieldValues(fieldName) {
				facet.addValue(fqValue, 0, true)
			}
		}

		facets = append(facets, facet)
	}
	return facets
}

func (r SearchResponse) newFacet(field string) FacetField {
	facet := FacetField{Field: field, Title: field}
	for _, def := range r.Params.Facets {
		if def.Field == facet.Field {
			// if the field is on the facets indicated on the params
			// (it should always be) grab the Title from it
			facet.Title = def.Title
			break
		}
	}
	return facet
}
