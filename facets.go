package solr

import (
	"strings"
)

type Facets []facetField

type facetField struct {
	Field  string
	Title  string
	Values []facetValue
}

// AddUrl and RemoveUrl are leaky abstraction since they are only
// needed in the user interface, but declaring them here simplify
// things a lot.
type facetValue struct {
	Text      string
	Count     int
	Active    bool   // true if we are filtering by this facet value
	AddUrl    string // URL to add this facet (set by the client)
	RemoveUrl string // URL to remove this facet (set by the client)
}

// Creates a new Facets object from a map. Notice that only facetFields
// are created in this case (i.e. no facetValues)
func NewFacets(definitions map[string]string) Facets {
	facets := Facets{}
	for key, value := range definitions {
		facets.Add(key, value)
	}
	return facets
}

func (facets *Facets) Add(field, title string) {
	facet := facetField{Field: field, Title: title}
	*facets = append(*facets, facet)
}

// Sets the AddUrl and RemoveUrl of the facet values for all the facets.
func (facets Facets) SetAddRemoveUrls(baseUrl string) {
	for _, facet := range facets {
		for i, value := range facet.Values {
			fqVal := "fq=" + facet.Field + "|" + value.Text + "&"
			facet.Values[i].RemoveUrl = strings.Replace(baseUrl, fqVal, "", 1)
			facet.Values[i].AddUrl = baseUrl + "&" + fqVal
		}
	}
}

func (ff *facetField) addValue(text string, count int, active bool) {
	value := facetValue{
		Text:   text,
		Count:  count,
		Active: active,
	}
	ff.Values = append(ff.Values, value)
}

func (facets Facets) toQueryString() string {
	qs := ""
	if len(facets) > 0 {
		qs += QsAdd("facet", "on")
		for _, f := range facets {
			qs += QsAdd("facet.field", f.Field)
			// The rest of the facet filters (mincount, limit, offset)
			// must be defined by the client.
		}
	}
	return qs
}
