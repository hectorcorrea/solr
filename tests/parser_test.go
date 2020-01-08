package tests

import (
	"fmt"
	"testing"

	"stash.uol.intranet/pdeng/solr-lib/solr"
	"stash.uol.intranet/pdeng/solrts"
)

func TestSelect(t *testing.T) {
	settings := solr.SettingsSolrCore{}
	var params map[string]string
	instance := solr.New("http://a1-kandango-q-pla1.host.intranet:8983", "pdeng_stats", false, true, settings, params, &solrts.TSDocumentParser{})
	facet := make(map[string]string)
	facet["facet.field"] = "metric"
	q := "{!parent which=parent_doc:true}"
	fl := "*,[child parentFilter=parent_doc:true]"

	searchParams := &solr.SearchParams{
		Q:      q,
		FL:     fl,
		Facets: facet,
		Rows:   1000,
	}

	res, _ := instance.Search(searchParams)
	fmt.Printf("%+v", res)
}
