package solr

import (
	"testing"
)

func TestDocument(t *testing.T) {
	d := Document{}
	d["single"] = "hello"
	d["multi-interface"] = []interface{}{"i1", "i2"}
	d["multi-string"] = []string{"s1", "s2"}

	if s := d.Value("single"); s != "hello" {
		t.Errorf("TestDocument unexpected value: %s", s)
	}

	if s := d.Values("multi-interface"); len(s) != 2 {
		t.Errorf("TestDocument unexpected value: %v", s)
	}

	if s := d.Values("multi-string"); len(s) != 2 {
		t.Errorf("TestDocument unexpected value: %v", s)
	}
}
