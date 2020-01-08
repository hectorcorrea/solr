package solr

import (
	"fmt"

	"github.com/buger/jsonparser"
)

const (
	rawMetric         string = "metric"
	rawStatus         string = "status"
	rawResponse       string = "response"
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
func (s *Instance) Decode(raw []byte) (Response, error) {
	res := Response{}
	var err error
	var i []interface{}
	if res.NumFound, res.Status, res.QTime, res.MaxScore, _ = s.parserNumbers(raw); err != nil {
		return Response{}, err
	}
	if res.Docs, err = s.DocumentParser.Parse(raw); err != nil {
		res.Docs = i
	}
	if res.Facets, _ = s.parserFacets(raw); err != nil {
		res.Facets = []facetField{}
	}

	if err != nil {
		return res, err
	}
	return res, err
}

func (s *Instance) parserFacets(raw []byte) ([]facetField, error) {
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

func (s *Instance) parserNumbers(raw []byte) (found, status, qtime int64, score float64, err error) {
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
