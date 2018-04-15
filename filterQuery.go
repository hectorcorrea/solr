package solr

import (
	"fmt"
	// "log"
	"net/url"
	"strings"
)

type filterQuery struct {
	Field string
	Value string
}

type filterQueries []filterQuery

// NewFilterQueries creates a new object from an array of values.
// values are the "fq=x|y" that came on the query string.
func newFilterQueries(values []string) filterQueries {
	fqs := filterQueries{}
	for _, value := range values {
		tokens := strings.Split(value, "|")
		if len(tokens) == 2 {
			fq := filterQuery{Field: tokens[0], Value: tokens[1]}
			fqs = append(fqs, fq)
		}
	}
	return fqs
}

func (fqs filterQueries) HasFieldValue(field, value string) bool {
	for _, fq := range fqs {
		if fq.Field == field && fq.Value == value {
			return true
		}
	}
	return false
}

func (fqs filterQueries) FieldValues(field string) []string {
	values := []string{}
	for _, fq := range fqs {
		if fq.Field == field {
			values = append(values, fq.Value)
		}
	}
	return values
}

func (fqs filterQueries) toQueryString() string {
	str := ""
	for _, fq := range fqs {
		str += fmt.Sprintf("fq=%s&", fq.toQueryString())
	}
	return str
}

func (fq filterQuery) toQueryString() string {
	// field:value, e.g. subject:"abc+xyz"
	return fmt.Sprintf("%s:%s", fq.Field, url.QueryEscape("\""+fq.Value+"\""))
}
