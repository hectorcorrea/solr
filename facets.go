package solr

import (
	"sort"
	"strconv"
	"strings"
)

// Facets is an array of FacetField definitions
type Facets []FacetField

// FacetField represents a single facet field definition.
type FacetField struct {
	Field  string       // Name of the field in Solr.
	Title  string       // Display title for the field.
	Order  int          // Order of this field on the Facets
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

// NewFacetsFromDefinitions creates a new Facets object from a map.
// Notice that only facetFields are created in this case (with a Field
// and Title, but no Values)
//
// If the Title for a definition is in the form "N|xxx" N is used as the
// order of the facet in the list.
func NewFacetsFromDefinitions(definitions map[string]string) Facets {
	facets := Facets{}
	for key, value := range definitions {
		tokens := strings.Split(value, "|")
		if len(tokens) < 2 {
			//no order indicated
			facets.add(key, value, 0)
		} else {
			order, _ := strconv.Atoi(tokens[0])
			facets.add(key, tokens[1], order)
		}
	}
	sort.Slice(facets, func(i, j int) bool { return facets[i].Order < facets[j].Order })
	return facets
}

func (facets *Facets) add(field, title string, order int) {
	facet := FacetField{Field: field, Title: title, Order: order}
	*facets = append(*facets, facet)
}

func (facets *Facets) ForField(field string) (FacetField, bool) {
	for _, facet := range *facets {
		if facet.Field == field {
			return facet, true
		}
	}
	return FacetField{}, false
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
			fqVal := "fq=" + facet.Field + "|" + value.Text + "&"
			facet.Values[i].RemoveUrl = strings.Replace(baseUrl, fqVal, "", 1)
			facet.Values[i].AddUrl = baseUrl + "&" + fqVal
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
