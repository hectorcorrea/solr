package solr

import (
	"net/url"
	"strings"
	"testing"
)

func TestGetParamsUrl(t *testing.T) {
	q := "id:123"
	fl := []string{"a", "b", "c"}
	options := map[string]string{"opt1": "val1"}
	params := NewGetParams(q, fl, options)
	qs := params.toSolrQueryString()
	if qs != "q=id%3A123&fl=a,b,c&opt1=val1&" {
		t.Errorf("Unexpected GetParams URL: %s", qs)
	}
}

func TestSearchParamsUrlEmpty(t *testing.T) {
	clientQs := url.Values{}
	options := map[string]string{}
	facets := map[string]string{}
	params := NewSearchParamsFromQs(clientQs, options, facets)
	qs := params.toSolrQueryString()
	if qs != "q=%2A&" {
		t.Errorf("Unexpected SearchParams URL: %s", qs)
	}
}

func TestSearchParamsUrl(t *testing.T) {
	clientQs := url.Values{
		"q":  []string{"title:\"one two\""},
		"fq": []string{"f1|v1", "f2|v2"},
	}
	options := map[string]string{"opt1": "val1"}
	facets := map[string]string{}
	params := NewSearchParamsFromQs(clientQs, options, facets)
	params.Fl = []string{"a", "b", "c"}
	qs := params.toSolrQueryString()
	if qs != "q=title%3A%22one+two%22&fl=a,b,c&fq=f1:%22v1%22&fq=f2:%22v2%22&opt1=val1&" {
		t.Errorf("Unexpected SearchParams URL: %s", qs)
	}

	facets = map[string]string{"faA": "xx", "faB": "yy"}
	params = NewSearchParamsFromQs(clientQs, options, facets)
	params.Fl = []string{"a", "b", "c"}
	qs = params.toSolrQueryString()
	if !strings.Contains(qs, "facet=on&") ||
		!strings.Contains(qs, "facet.field=faA&") ||
		!strings.Contains(qs, "facet.field=faB&") {
		t.Errorf("Unexpected SearchParams (facets) URL: %s", qs)
	}
}
