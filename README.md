# Solr
A Solr client written in Go.

This client is geared towards supporting a web user interface that queries and filters via facets. This package provides a lot of the functionality to know what facets the user is filtering on and allow the user interface to add and remove more facets to the existing search.

Project https://github.com/hectorcorrea/solrdora is an example of a complete web application using this package.


## Source code

* `solr.go`: the main entry point, functions `Get()` and `Search()` are defined here.
* `getParams.go`: parameter for the `Search()` function.
* `searchParams.go`: parameter for the `Search()` function.
* `searchResponse.go`: the object used to represent the results of a `Search()`
* `filterQueries.go`: represent the `fq` values passed to Solr.


## Examples of use (basic)

Search for documents
```
q := "title:\"one two three\""

solr := solr.New("http://localhost:8983/solr/your-core")
results, err := solr.SearchText(q)

log.Printf("Documents found: %d", results.NumFound)
for i, doc := results.Documents {
  log.Printf("%d %v", i, doc)
}
```

Get one Solr document by ID

```
q := "id:123"
fl := []string{}
options := map[string]string{}
params := NewGetParams(q, fl, options)

solr := solr.New("http://localhost:8983/solr/your-core")
doc, err := solr.Get(params)

log.Printf("%v", doc)
```

## More examples
Search for documents customizing list of fields to retrieve,
facets, and other parameters.
```
# In a web app qs will be req.URL.Query() where req is
# the *http.Request that you get in your HTTP handler.
qs := url.Values{
  "q":  []string{"title:\"one two\""},
  "fq": []string{"subject|Geography"},
}

options := map[string]string{
  "defType": "edismax",
}

facets := map[string]string{
  "publisher": "Publisher name",
  "subject_str": "Subject",
}

params := NewSearchParamsFromQs(qs, options, facets)
params.Fl = []string{"id", "title", "authorsAll", "_version_"}

solr := solr.New("http://localhost:8983/solr/your-core")
results, err := solr.Search(params)

log.Printf("Documents found: %d", results.NumFound)
for i, doc := results.Documents {
  log.Printf("%d %v", i, doc)
}
```

A full-blow example using this library can be found in the
[SolrDora](https://github.com/hectorcorrea/solrdora) repo.

The GoDoc documentation for this package can be found [here](https://godoc.org/github.com/hectorcorrea/solr)
