package solr

import (
	"fmt"
	"reflect"
	"strings"
)

type Document struct {
	Data       map[string]interface{}
	Highlights map[string][]string
}

func NewDocument() Document {
	data := map[string]interface{}{}
	hl := map[string][]string{}
	return Document{Data: data, Highlights: hl}
}

func NewDocumentFromSolrDoc(data documentRaw) Document {
	hl := map[string][]string{}
	return Document{Data: data, Highlights: hl}
}

func NewDocumentFromSolrResponse(raw responseRaw) []Document {
	docs := []Document{}
	for _, rawDoc := range raw.Data.Documents {
		// Create the document...
		doc := NewDocumentFromSolrDoc(rawDoc)

		// ...and attach its highlight information from the Solr response
		for field, values := range raw.Highlighting[doc.Id()] {
			doc.Highlights[field] = values
		}

		docs = append(docs, doc)
	}
	return docs
}

// Returns the value in a field. Concatenates multi-value
// fields into a single string.
func (d Document) Value(fieldName string) string {
	values := d.Values(fieldName)
	return strings.Join(values, " ")
}

// Returns all the values in a multi-value field
// (mimics an array of one if the field is single value)
func (d Document) Values(fieldName string) []string {
	var values []string
	dynValue := reflect.ValueOf(d.Data[fieldName])
	kind := dynValue.Kind()
	if kind == reflect.Invalid {
		return values
	}
	if kind == reflect.Slice {
		for i := 0; i < dynValue.Len(); i++ {
			strValue := fmt.Sprintf("%s", dynValue.Index(i))
			values = append(values, strValue)
		}
		return values
	}
	strValue := fmt.Sprintf("%s", dynValue)
	values = append(values, strValue)
	return values
}

func (d Document) ValueFloat(fieldName string) float64 {
	value, ok := d.Data[fieldName].(float64)
	if ok {
		return value
	}
	return 0.0
}

func (d Document) Id() string {
	return d.Value("id")
}

func (d Document) HighlightsFor(field string) []string {
	return d.Highlights[field]
}

func (d Document) HighlightFor(field string) string {
	values := d.Highlights[field]
	if len(values) > 0 {
		return strings.Join(values, " ")
	}
	return ""
}

func (d Document) IsHighlighted(field string) bool {
	return len(d.Highlights[field]) > 0
}
