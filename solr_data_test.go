package solr

import (
	"net/url"
	"testing"
)

// These tests depends an instance of Solr
// running at http://localhost:8983/solr/bibdata
const solrCoreUrl = "http://localhost:8983/solr/bibdata"

// const solrCoreUrl = "http://localhost:8081/solr/blacklight-core"

func xTestCount(t *testing.T) {
	solr := New(solrCoreUrl, true)
	_, err := solr.Count()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

func xTestGetData(t *testing.T) {
	q := "id:00009565"
	fl := []string{}
	options := map[string]string{}
	params := NewGetParams(q, fl, options)
	solr := New(solrCoreUrl, false)
	_, err := solr.Get(params)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

func xTestHighlightsData(t *testing.T) {
	qs := url.Values{
		"q": []string{"george"},
	}
	options := map[string]string{
		"defType": "edismax",
		"qf":      "authorsAll title",
		"hl":      "on",
	}
	facets := map[string]string{}
	params := NewSearchParamsFromQs(qs, options, facets)
	params.Fl = []string{"id", "authorsAll", "title"}

	solr := New(solrCoreUrl, false)
	results, err := solr.Search(params)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	doc := results.Documents[0]
	if !doc.IsHighlighted("title") {
		t.Errorf("Expected title to be highlighted")
	}

	if doc.IsHighlighted("authorx") {
		t.Errorf("Unexpected field was highlighted")
	}
}

func xTestPostData(t *testing.T) {
	doc := Document{
		Data: map[string]interface{}{
			"id":           "123",
			"author":       "ada lovelace",
			"authorsOther": []string{"a", "b"},
		},
	}
	solr := New(solrCoreUrl, false)
	err := solr.PostDoc(doc)
	if err != nil {
		t.Errorf("PostDoc error: %s", err)
	}

	data := map[string]interface{}{
		"id":     "456",
		"author": "grace hopper",
	}
	err = solr.PostOne(data)
	if err != nil {
		t.Errorf("PostOne error: %s", err)
	}
}

func xTestDeleteData(t *testing.T) {
	ids := []string{"00000092", "00000093"}
	solr := New(solrCoreUrl, false)
	err := solr.Delete(ids)
	if err != nil {
		t.Errorf("PostDoc error: %s", err)
	}
}
