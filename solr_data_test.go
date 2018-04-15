package solr

// These tests depends an instance of Solr
// running at http://localhost:8983/solr/bibdata

import (
	"net/url"
	"testing"
)

func TestGetData(t *testing.T) {
	q := "id:00009565"
	fl := []string{}
	options := map[string]string{}
	params := NewGetParams(q, fl, options)
	solr := New("http://localhost:8983/solr/bibdata", false)
	_, err := solr.Get(params)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

// Depends on Solr running at http://localhost:8983/solr/bibdata
func TestHighlightsData(t *testing.T) {
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

	solr := New("http://localhost:8983/solr/bibdata", false)
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

func TestPostData(t *testing.T) {
	doc1 := Document{
		Data: map[string]interface{}{
			"id":           "123",
			"author":       "ada lovelace",
			"authorsOther": []string{"a", "b", "c"},
		},
	}

	doc2 := Document{
		Data: map[string]interface{}{
			"id":           "456",
			"author":       "grace hopper",
			"authorsOther": []string{"x", "y"},
		},
	}

	solr := New("http://localhost:8983/solr/bibdata", false)
	docs := []Document{doc1, doc2}
	err := solr.PostDocs(docs)
	if err != nil {
		t.Errorf("PostDocs error: %s", err)
	}

	doc3 := Document{
		Data: map[string]interface{}{
			"id":           "789",
			"author":       "karla loya",
			"authorsOther": []string{"k", "l"},
		},
	}
	err = solr.PostDoc(doc3)
	if err != nil {
		t.Errorf("PostDoc error: %s", err)
	}

	data := map[string]interface{}{
		"id":     "000",
		"author": "mac loya",
	}
	err = solr.PostOne(data)
	if err != nil {
		t.Errorf("PostOne error: %s", err)
	}
}
