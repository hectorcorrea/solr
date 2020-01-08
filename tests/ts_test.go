package tests

import (
	"fmt"
	"testing"

	"stash.uol.intranet/pdeng/solr-lib/solr"
	"stash.uol.intranet/pdeng/solrts"
)

func BenchmarkSolrLib(b *testing.B) {
	settings := solr.SettingsSolrCore{}
	var params map[string]string
	instance := solr.New("http://localhost:8080", "produtos_digitais", false, false, settings, params, &solrts.TSDocumentParser{})
	searchParams := &solr.SearchParams{}
	res, err := instance.Search(searchParams)
	if err != nil {
		b.Error(err.Error())
	}
	fmt.Printf("%+v", res)
}
