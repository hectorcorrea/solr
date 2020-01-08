package manager

import (
	"fmt"

	"github.com/buger/jsonparser"
	"stash.uol.intranet/pdeng/solr-lib/parser"
)

//Response - response data from solr instance
type Response struct {
	Status   int64        `json:"status,omitempty"`
	QTime    int64        `json:"Qtime,omitempty"`
	NumFound int64        `json:"numFound,omitempty"`
	MaxScore float64      `json:"maxScore,omitempty"`
	Docs     interface{}  `json:"Docs,omitempty"`
	Facets   []facetField `json:"Facets,omitempty"`
}

type facetField struct {
	name string
	list []facetValue
}
type facetValue struct {
	name  string
	value int64
}

//Manager - manager for custom implement
type Manager struct {
	documentParser parser.DocumentParser
}

//New - creare a new manager
func New(documentParser parser.DocumentParser) *Manager {
	return &Manager{
		documentParser: documentParser,
	}
}

const (
	rawMetric         string = "metric"
	rawStatus         string = "status"
	rawResponse       string = "response"
	rawDoc            string = "docs"
	rawChildDocs      string = "_childDocuments_"
	rawFacetsCount    string = "facet_counts"
	rawFacetFields    string = "facet_fields"
	rawNumFound       string = "numFound"
	rawResponseHeader string = "responseHeader"
	rawQtime          string = "QTime"
	rawMaxScore       string = "maxScore"
	rawArrayString    string = "[0]"
)

//Decode - decode a raw byte from solr instance and return a formated response
func (m *Manager) Decode(raw []byte) (Response, error) {
	res := Response{}
	var err error
	var i []interface{}
	if res.NumFound, res.Status, res.QTime, res.MaxScore, _ = m.parserNumbers(raw); err != nil {
		return Response{}, err
	}
	if res.Docs, err = m.documentParser.Parse(raw); err != nil {
		res.Docs = i
	}
	if res.Facets, _ = m.parserFacets(raw); err != nil {
		res.Facets = []facetField{}
	}

	if err != nil {
		return res, err
	}
	return res, err
}

func (m *Manager) parserFacets(raw []byte) ([]facetField, error) {
	facetValueArray := []facetValue{}
	facetValues := facetValue{}
	facetFields := facetField{}
	var facetFieldsArray []facetField
	var err error
	var k string
	var v int64
	err = jsonparser.ObjectEach(raw, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {

		_, err = jsonparser.ArrayEach(value, func(tvalue []byte, dataType jsonparser.ValueType, offset int, err error) {

			if err != nil {
				return
			}
			switch dataType {
			case jsonparser.String:
				k = fmt.Sprintf("%s", tvalue)
			case jsonparser.Number:
				v, err = jsonparser.GetInt(tvalue)
				facetValues.name = k
				facetValues.value = v
				facetValueArray = append(facetValueArray, facetValues)
			}
		})
		if err != nil {
			return err
		}

		facetFields.name = string(key)
		facetFields.list = facetValueArray
		facetFieldsArray = append(facetFieldsArray, facetFields)
		return nil
	}, rawFacetsCount, rawFacetFields)
	if err != nil {
		return nil, err
	}

	return facetFieldsArray, nil
}

func (m *Manager) parserNumbers(raw []byte) (found, status, qtime int64, score float64, err error) {
	raw = jsonparser.Delete(raw, rawResponseHeader)
	if found, err = jsonparser.GetInt(raw, rawResponse, rawNumFound); err != nil {
		return found, status, qtime, score, err
	}
	if status, err = jsonparser.GetInt(raw, rawResponseHeader, rawStatus); err != nil {
		return found, status, qtime, score, err
	}
	if qtime, err = jsonparser.GetInt(raw, rawResponseHeader, rawQtime); err != nil {
		return found, status, qtime, score, err
	}
	if score, err = jsonparser.GetFloat(raw, rawResponse, rawMaxScore); err != nil {
		return found, status, qtime, score, err
	}

	return found, status, qtime, score, err
}
