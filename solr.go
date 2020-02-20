// Package solr provides functions to connect to Solr and make request
// for getting individual documents, executing searches, updating, and
// deleteing documents.
//
// This package is geared towards supporting a web user interface that
// queries and filters Solr via facets. As such it provides functionality
// to handle the typical request-response workflow of a web application.
// For example SearchResponse provides URLs to re-execute a search and
// handle pagination, likewise the Facets returned in a SearchResponse
// include URLs to add or remove a filter for a given facet field/value.
//
// Most basic search usage:
//
// 	s := solr.New("http://localhost/solr/some-core", false)
// 	results := s.SearchText("hello")
// 	log.Printf("%v", results)
//
// Search with options:
// 	s := solr.New("http://localhost/solr/some-core", false)
//	qs := url.Values{
//		"q": []string{"title:\"one two\""},
//	}
// 	options := map[string]interface{}{
//		"defType": "edismax",
// 	}
//	facets := map[string]string{
//		"title_str" : "Title",
//	}
//	params := NewSearchParams(qs, options, facets)
//	results := s.Search(params)
//	log.Printf("%v", results)
//
// Search with options from a query string (req is *http.Request
// from a web handler)
//
// 	s := solr.New("http://localhost/solr/some-core", false)
// 	options := map[string]interface{}{}
//	facets := map[string]string{}
//	params := NewSearchParams(req.URL.Query(), options, facets)
//	results := s.Search(params)
//	log.Printf("%v", results)
//
package solr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// The main class to drive interaction with Solr.
type Solr struct {
	CoreUrl string
	Verbose bool
}

// Creates a new instance of Solr.
// When verbose = true it will log.Printf() the HTTP requests to Solr.
func New(coreUrl string, verbose bool) Solr {
	return Solr{CoreUrl: coreUrl, Verbose: verbose}
}

func (s Solr) Count() (int, error) {
	options := map[string]string{"rows": "0", "defType": "edismax", "wt": "json"}
	facets := map[string]string{}
	params := NewSearchParams("*", options, facets)
	r, err := s.Search(params)
	return r.NumFound, err
}

// Get fetches a single document from Solr.
func (s Solr) Get(params GetParams) (Document, error) {
	url := s.CoreUrl + "/select?" + params.toSolrQueryString()
	raw, err := s.httpGet(url)
	if err != nil {
		return Document{}, err
	}

	count := len(raw.Data.Documents)
	if count == 0 {
		return Document{}, nil
	} else if count > 1 {
		msg := fmt.Sprintf("More than one document was found (Q=%s)", params.Q)
		return Document{}, errors.New(msg)
	}
	return newDocumentFromSolrDoc(raw.Data.Documents[0]), err
}

// Issues a search with the values indicated in the paramers.
func (s Solr) Search(params SearchParams) (SearchResponse, error) {
	url := s.CoreUrl + "/select?" + params.toSolrQueryString()
	raw, err := s.httpGet(url)
	if err != nil {
		return SearchResponse{}, err
	}
	return newSearchResponse(params, raw), err
}

// Updates a single document in Solr with the data in the
// document provided.
func (s Solr) PostDoc(doc Document) error {
	docs := []Document{doc}
	return s.PostDocs(docs)
}

// Updates an array of documents in Solr.
func (s Solr) PostDocs(docs []Document) error {
	// Extract the data from the documents
	// (i.e. only the data, without the highlight properties)
	data := []map[string]interface{}{}
	for _, doc := range docs {
		data = append(data, doc.Data)
	}
	return s.Post(data)
}

// Updates a single document in Solr. Uses plain Go map[string]interface{}
// object rather than a Document object. The map key is represents
// the field name and the map value the field value.
func (s Solr) PostOne(datum map[string]interface{}) error {
	data := []map[string]interface{}{datum}
	return s.Post(data)
}

// Updates an array of documents in Solr. Uses an array of
// plain Go map[string]interface{} object rather than an
// array of Document objects. The map key is represents
// the field name and the map value the field value.
func (s Solr) Post(data []map[string]interface{}) error {
	contentType := "application/json"
	params := "wt=json&commit=true"
	url := s.CoreUrl + "/update?" + params
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r, err := s.httpPost(url, contentType, string(bytes))
	if err != nil {
		return err
	}

	var response responseRaw
	err = json.Unmarshal([]byte(r), &response)
	if err != nil {
		return err
	}

	if response.Header.Status != 0 {
		errorMsg := fmt.Sprintf("Solr returned status %d", response.Header.Status)
		return errors.New(errorMsg)
	}

	return nil
}

// Deleteds from Solr the documents with the IDs indicated.
func (s Solr) Delete(ids []string) error {
	// notice that the request body (contentType) is in XML
	// but the response (wt) is in JSON
	contentType := "text/xml"
	params := "wt=json&commit=true"
	url := s.CoreUrl + "/update?" + params

	payload := "<delete>\r\n"
	for _, id := range ids {
		payload += "\t<id>" + id + "</id>\r\n"
	}
	payload += "</delete>"

	r, err := s.httpPost(url, contentType, payload)
	if err != nil {
		return err
	}

	var response responseRaw
	err = json.Unmarshal([]byte(r), &response)
	if err != nil {
		return err
	}

	if response.Header.Status != 0 {
		errorMsg := fmt.Sprintf("Solr returned status %d", response.Header.Status)
		return errors.New(errorMsg)
	}

	return nil
}

// Issues a search for the text indicated. Uses the server's default
// values for all other Solr parameters.
func (s Solr) SearchText(text string) (SearchResponse, error) {
	options := map[string]string{}
	facets := map[string]string{}
	params := NewSearchParams(text, options, facets)
	return s.Search(params)
}

func (s Solr) httpGet(url string) (responseRaw, error) {
	s.log("Solr HTTP GET", url)
	r, err := http.Get(url)
	if err != nil {
		return responseRaw{}, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return responseRaw{}, err
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		msg := fmt.Sprintf("HTTP Status: %s. ", r.Status)
		if len(body) > 0 {
			msg += fmt.Sprintf("Body: %s", body)
		}
		return responseRaw{}, errors.New(msg)
	}

	response, err := NewResponseRaw([]byte(body))
	if err == nil {
		// HTTP request was successful but Solr reported an error.
		if response.Error.Trace != "" {
			msg := fmt.Sprintf("Solr Error. %#v", response.Error)
			err = errors.New(msg)
		}
	} else {
		if len(r.Header["Content-Type"]) > 0 {
			// Perhaps the response was not in JSON
			// (e.g. if Solr returns XML by default)
			msg := fmt.Sprintf("%s. Solr's response Content-Type: %s", err, r.Header["Content-Type"])
			err = errors.New(msg)
		}
	}
	return response, err
}

func (s Solr) httpPost(url, contentType, body string) (string, error) {
	s.log("Solr HTTP POST", url)
	payload := bytes.NewBufferString(body)
	r, err := http.Post(url, contentType, payload)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	respStr, err := ioutil.ReadAll(r.Body)
	return string(respStr), nil
}

func (s Solr) log(msg1, msg2 string) {
	if s.Verbose {
		log.Printf("%s: %s", msg1, msg2)
	}
}
