// Package solr provides functions to connect to Solr and make request
// for getting individual documents, executing searches, updating, and
// deleteing documents.
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

type Solr struct {
	CoreUrl string
	Verbose bool
}

func New(coreUrl string, verbose bool) Solr {
	return Solr{CoreUrl: coreUrl, Verbose: verbose}
}

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
	return NewDocumentFromSolrDoc(raw.Data.Documents[0]), err
}

// Issues a search with the values indicated in the paramers
func (s Solr) Search(params SearchParams) (SearchResponse, error) {
	url := s.CoreUrl + "/select?" + params.toSolrQueryString()
	raw, err := s.httpGet(url)
	if err != nil {
		return SearchResponse{}, err
	}
	return NewSearchResponse(params, raw), err
}

func (s Solr) PostDoc(doc Document) error {
	docs := []Document{doc}
	return s.PostDocs(docs)
}

func (s Solr) PostDocs(docs []Document) error {
	// Extract the data from the documents
	// (i.e. only the data, without the highlight properties)
	data := []map[string]interface{}{}
	for _, doc := range docs {
		data = append(data, doc.Data)
	}
	return s.Post(data)
}

func (s Solr) PostOne(datum map[string]interface{}) error {
	data := []map[string]interface{}{datum}
	return s.Post(data)
}

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

// Issues a search for the text indicated and using only
// Solr default values
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

	// log.Printf("Body: %s", body)

	var response responseRaw
	err = json.Unmarshal([]byte(body), &response)
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
