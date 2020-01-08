package solr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
	"stash.uol.intranet/pdeng/solr-lib/parser"
)

// Instance The main class to instance solr
type Instance struct {
	CoreURL           string //CoreURL - ex: http://localhost:8983
	InstaceName       string //InstanceName - Name of Collection/Core
	SolrCore          bool   //SolrCore - set true if use Solr Core instance
	BlockJoinFaceting bool
	CoreConfig        SettingsSolrCore  //CoreConfig - Configs required for admin solr core
	CloudParams       map[string]string //CloudParams - Params used for admin solr cloud
	DocumentParser    parser.DocumentParser
	listCollectionURL string
	newSolrCoreURL    string
	newSolrCloudURL   string
	deleteURL         string
}

// SearchParams - Params for solr queries
type SearchParams struct {
	Q             string
	FL            string
	FilterQueries map[string]string
	Sort          string
	Facets        map[string]string
	Rows          int
	Start         int
}

//ResponseFromSolr response from solr instance
type ResponseFromSolr struct {
	Status   int64
	QTime    int64
	NumFound int
	Response ResponseDocs
}

// CreateResponse - Solr Response for parser message errors
type CreateResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Error struct {
		Metadata []string `json:"metadata"`
		Msg      string   `json:"msg"`
		Trace    string   `json:"trace"`
		Code     int      `json:"code"`
	} `json:"error"`
}

// ResponseDocs - Solr response Docs
type ResponseDocs struct {
	Docs      map[string]string
	ChildDocs map[string]string
}

// SettingsSolrCore - Configurations for create a new solr core for more information visit https://lucene.apache.org/solr/guide/6_6/coreadmin-api.html
type SettingsSolrCore struct {
	CoreName    string //CoreName - The name of the new core. Same as "name" on the <core> element.
	InstanceDir string //InstanceDir - The directory where files for this SolrCore should be stored. Same as instanceDir on the <core> element.
	Config      string //Config - Name of the config file (i.e., solrconfig.xml) relative to instanceDir
	Schema      string //Schema - Name of the schema file to use for the core.
	DataDir     string //DataDir - Name of the data directory relative to instanceDir
}

//ResponseData array of solr response
type ResponseData []byte

const (
	rawParserStatus      string = "status"
	rawParserCollections string = "collections"
)

//New Create a new instance os Solr
func New(coreURL, instanceName string, solrCore, BlockJoinFaceting bool, CoreConfig SettingsSolrCore, CloudParams map[string]string, documentParser parser.DocumentParser) Instance {
	var listurl, pr, newcoreurl, newcloudurl, deleteurl string
	if solrCore {
		deleteurl = coreURL + "/solr/admin/cores?action=UNLOAD&core=" + instanceName + "&deleteInstanceDir=true"
		listurl = coreURL + "/solr/admin/cores?action=STATUS&wt=json"
		newcoreurl = coreURL + "/solr/admin/cores?action=CREATE&name=" + CoreConfig.CoreName + "&instanceDir=" + CoreConfig.InstanceDir + "&config=" + CoreConfig.InstanceDir + "&schema=" + CoreConfig.Schema + "&dataDir=" + CoreConfig.DataDir
	} else {
		listurl = coreURL + "/solr/admin/collections?action=LIST&indexInfo=false&wt=json"
		if CloudParams != nil {
			for k, v := range CloudParams {
				pr += fmt.Sprintf("&%s=%s", url.QueryEscape(k), url.QueryEscape(v))
			}
			newcloudurl = coreURL + "/solr/admin/collections?action=CREATE&name=" + instanceName + pr

		} else {
			deleteurl = coreURL + "/solr/admin/collections?action=DELETE&name=" + instanceName
			newcloudurl = coreURL + "/solr/admin/collections?action=CREATE&name=" + instanceName
		}
	}
	return Instance{CoreURL: coreURL, InstaceName: instanceName, CoreConfig: CoreConfig, SolrCore: solrCore, CloudParams: CloudParams, DocumentParser: documentParser, listCollectionURL: listurl, newSolrCoreURL: newcoreurl, newSolrCloudURL: newcloudurl, deleteURL: deleteurl}
}

// NewSolrCore create a new solr core
func (s *Instance) newSolrCore(newcoreurl string) error {
	if newcoreurl == "" {
		return errors.New("new core url not defined")
	}
	var response *CreateResponse
	res, err := http.Get(newcoreurl)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &response)
	if err != nil || res.StatusCode != 200 {
		return errors.New(response.Error.Msg)
	}
	return nil
}

func (s *Instance) newSolrCloud(url string) error {
	if url == "" {
		return errors.New("new cloud url not defined")
	}
	var (
		response *CreateResponse
	)
	res, err := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(body), &response)
	if err != nil || res.StatusCode != 200 {
		return errors.New(response.Error.Msg)
	}
	return nil
}

// List - List cores/collections
func (s *Instance) List() ([]interface{}, error) {
	var cores []interface{}
	raw, err := s.httpGet(s.listCollectionURL)
	if err != nil {
		return nil, err
	}
	i, err := s.listParser(raw)
	if err != nil {
		fmt.Println(err)
	}
	cores = i
	return cores, nil
}

func (s *Instance) listParser(raw []byte) ([]interface{}, error) {
	var list []interface{}
	switch s.SolrCore {
	case true:
		err := jsonparser.ObjectEach(raw, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
			name, _, _, err := jsonparser.Get(value, "name")
			if err != nil {
				return err
			}
			list = append(list, string(name))
			return nil
		}, rawParserStatus)
		if err != nil {
			return nil, err
		}
	default:
		_, err := jsonparser.ArrayEach(raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			list = append(list, string(value))
		}, rawParserCollections)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

// Delete - delete a core/collection
func (s *Instance) Delete() error {
	_, err := s.httpGet(s.deleteURL)
	if err != nil {
		return err
	}
	return nil
}

// Create - create cores/collections
func (s *Instance) Create() error {
	if s.SolrCore && s.InstaceName != "" {
		err := s.newSolrCore(s.newSolrCoreURL)
		if err != nil {
			return err
		}
	} else {
		err := s.newSolrCloud(s.newSolrCloudURL)
		if err != nil {
			return err
		}
	}
	return nil
}

//Search A basic search in solr
func (s *Instance) Search(params *SearchParams) (Response, error) {
	var url string
	var res Response
	if s.BlockJoinFaceting {
		url = s.CoreURL + "/solr/" + s.InstaceName + "/bjqfacet?" + params.toQueryString()
	} else {
		url = s.CoreURL + "/solr/" + s.InstaceName + "/select?" + params.toQueryString()
	}
	raw, err := s.httpGet(url)
	if err != nil {
		return res, err
	}

	res, err = s.Decode(raw)
	if err != nil {
		return Response{}, err
	}
	return res, nil
}

func (params SearchParams) toQueryString() string {

	qs := ""
	if params.Sort != "" {
		qs += fmt.Sprintf("&sort=%s", url.QueryEscape(params.Sort))
	}
	if params.Q != "" {
		qs += fmt.Sprintf("&q=%s", url.QueryEscape(params.Q))
	}
	if params.FL != "" {
		qs += fmt.Sprintf("&fl=%s", url.QueryEscape(params.FL))
	}
	if params.FilterQueries != nil {
		for k, v := range params.FilterQueries {
			qs += fmt.Sprintf("&fq=%s:%s", url.QueryEscape(k), url.QueryEscape(v))
		}
	}
	if params.Facets != nil {
		for k, v := range params.Facets {
			qs += fmt.Sprintf("&facet=true&%s=%s", url.QueryEscape(k), url.QueryEscape(v))
		}
	}
	if params.Start != 0 {
		qs += fmt.Sprintf("&start=%d", params.Start)
	}
	if params.Rows != 0 {
		qs += fmt.Sprintf("&rows=%d", params.Rows)
	}
	return qs
}

// Post - post json on solr, if the postParams is passed it will be adding in the request. For deleting items you can use the post function using the json format: https://lucene.apache.org/solr/guide/6_6/uploading-data-with-index-handlers.html#UploadingDatawithIndexHandlers-SendingJSONUpdateCommands
func (s *Instance) Post(payload []byte, postParams map[string]string) error {
	contentType := "application/json"
	params := ""
	if postParams != nil {
		for k, v := range postParams {
			params += fmt.Sprintf("&%s=%s", url.QueryEscape(k), url.QueryEscape(v))
		}
	}
	url := s.CoreURL + "/solr/" + s.InstaceName + "/update?" + params + "&wt=json"

	r, err := s.httpPost(url, contentType, string(payload))
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

func (s *Instance) httpPost(url, contentType, body string) (string, error) {
	payload := bytes.NewBufferString(body)
	res, err := http.Post(url, contentType, payload)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resSTR, err := ioutil.ReadAll(res.Body)
	return string(resSTR), nil
}
func (s *Instance) httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		msg := fmt.Sprintf("HTTP Status: %s. ", res.Status)
		if len(body) > 0 {
			msg += fmt.Sprintf("Body: %s", body)
		}
		return nil, fmt.Errorf("HTTP Status: %s", msg)
	}
	return body, nil
}
