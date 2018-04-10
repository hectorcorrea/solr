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

// Returns the value in a single-value field
func (d Document) Value(fieldName string) string {
	// Casting to string would have been cleaner but it _only_ works for strings.
	// Casting to interface{} allows us to fetch the value even if it is not
	// a string (e.g a float). The downside is that fmt.Sprintf() returns a
	// funny value for non-strings, but at least we fetch the value.
	value, ok := d.Data[fieldName].(interface{})
	if ok {
		return fmt.Sprintf("%s", value)
	}
	return ""
}

// Returns all the values in a multi-value field
func (d Document) Values(fieldName string) []string {
	var values []string
	dynValue := reflect.ValueOf(d.Data[fieldName])
	if dynValue.Kind() == reflect.Slice {
		for i := 0; i < dynValue.Len(); i++ {
			strValue := fmt.Sprintf("%s", dynValue.Index(i))
			values = append(values, strValue)
		}
	}
	return values
}

// Returns the first value in a multi-value field
func (d Document) FirstValue(fieldName string) string {
	values := d.Values(fieldName)
	if len(values) == 0 {
		return ""
	}
	return values[0]
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
