package solr

import (
	"net/url"
	"strings"
)

// An array of FacetField definitions
type Facets []FacetField

// A single facet field definition.
type FacetField struct {
	Field  string       // Name of the field in Solr.
	Title  string       // Display title for the field.
	Values []FacetValue // Values returned by Solr for this field.
}

// AddUrl and RemoveUrl are leaky abstraction since they are only
// needed in the user interface, but declaring them here simplify
// things a lot upstream.
type FacetValue struct {
	Text      string // Value returned by Solr for this field.
	Count     int    // Number of documents that matched this field/value.
	Active    bool   // true if we are filtering by this facet value
	AddUrl    string // URL to filter by this value. See SetAddRemoveUrls()
	RemoveUrl string // URL to remove filter by this value. See SetAddRemoveUrls()
}

// Creates a new Facets object from a map. Notice that only facetFields
// are created in this case (with a Field and Title, but no Values)
func newFacets(definitions map[string]string) Facets {
	facets := Facets{}
	for key, value := range definitions {
		facets.add(key, value)
	}
	return facets
}

func (facets *Facets) add(field, title string) {
	facet := FacetField{Field: field, Title: title}
	*facets = append(*facets, facet)
}

// Sets the internal AddUrl and RemoveUrl of the facet values for
// all the facets using the provided baseUrl.
//
// AddUrl is a URL that can be used in NewSearchParamsFromQs() to create
// a search filtering by the field/value of the facet.
//
// RemoveUrl is a URL that can be used in NewSearchParamsFromQs() to create
// a search not filtering by this field/value in the search.
func (facets Facets) SetAddRemoveUrls(baseUrl string) {
	for _, facet := range facets {
		for i, value := range facet.Values {
			fqValRaw := "fq=" + facet.Field + "|" + value.Text + "&"
			facet.Values[i].RemoveUrl = strings.Replace(baseUrl, fqValRaw, "", 1)
			fqValEncoded := "fq=" + facet.Field + "|" + url.QueryEscape(value.Text) + "&"
			facet.Values[i].AddUrl = baseUrl + "&" + fqValEncoded
		}
	}
}

func (ff *FacetField) addValue(text string, count int, active bool) {
	value := FacetValue{
		Text:   text,
		Count:  count,
		Active: active,
	}
	ff.Values = append(ff.Values, value)
}

func (facets Facets) toQueryString() string {
	qs := ""
	if len(facets) > 0 {
		qs += qsAdd("facet", "on")
		for _, f := range facets {
			qs += qsAdd("facet.field", f.Field)
			// The rest of the facet filters (mincount, limit, offset)
			// must be defined by the client.
		}
	}
	return qs
}
