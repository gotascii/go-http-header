package header

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewFromStruct_types(t *testing.T) {
	str := "string"
	strPtr := &str

	tests := []struct {
		in   interface{}
		want http.Header
	}{
		{
			// basic primitives
			struct {
				A string
				B int
				C uint
				D float32
				E bool
			}{},
			http.Header{
				"A": {""},
				"B": {"0"},
				"C": {"0"},
				"D": {"0"},
				"E": {"false"},
			},
		},
		{
			// pointers
			struct {
				A *string
				B *int
				C **string
			}{A: strPtr, C: &strPtr},
			http.Header{
				"A": {str},
				"B": {""},
				"C": {str},
			},
		},
		{
			// other types
			struct {
				A time.Time
				B time.Time `header:",unix"`
				C bool      `header:",int"`
				D bool      `header:",int"`
			}{
				A: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				B: time.Date(2000, 1, 1, 12, 34, 56, 0, time.UTC),
				C: true,
				D: false,
			},
			http.Header{
				"A": {"2000-01-01T12:34:56Z"},
				"B": {"946730096"},
				"C": {"1"},
				"D": {"0"},
			},
		},
		{
			nil,
			http.Header{},
		},
	}

	for i, tt := range tests {
		v, err := NewFromStruct(tt.in)
		if err != nil {
			t.Errorf("%d. NewFromStruct(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. NewFromStruct(%q) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestNewFromStruct_omitEmpty(t *testing.T) {
	str := ""
	s := struct {
		a string
		A string
		B string  `header:",omitempty"`
		C string  `header:"-"`
		D string  `header:"omitempty"` // actually named omitempty, not an option
		E *string `header:",omitempty"`
	}{E: &str}

	v, err := NewFromStruct(s)
	if err != nil {
		t.Errorf("NewFromStruct(%q) returned error: %v", s, err)
	}

	want := http.Header{
		"A":         {""},
		"Omitempty": {""},
		"E":         {""}, // E is included because the pointer is not empty, even though the string being pointed to is
	}
	if !reflect.DeepEqual(want, v) {
		t.Errorf("NewFromStruct(%q) returned %v, want %v", s, v, want)
	}
}

type A struct {
	B
}

type B struct {
	C string
}

type D struct {
	B
	C string
}

type e struct {
	B
	C string
}

type F struct {
	e
}

func TestNewFromStruct_embeddedStructs(t *testing.T) {
	tests := []struct {
		in   interface{}
		want http.Header
	}{
		{
			A{B{C: "foo"}},
			http.Header{"C": {"foo"}},
		},
	}

	for i, tt := range tests {
		v, err := NewFromStruct(tt.in)
		if err != nil {
			t.Errorf("%d. NewFromStruct(%q) returned error: %v", i, tt.in, err)
		}

		if !reflect.DeepEqual(tt.want, v) {
			t.Errorf("%d. NewFromStruct(%q) returned %v, want %v", i, tt.in, v, tt.want)
		}
	}
}

func TestNewFromStruct_invalidInput(t *testing.T) {
	_, err := NewFromStruct("")
	if err == nil {
		t.Errorf("expected NewFromStruct() to return an error on invalid input")
	}
}

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
		{"field", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}
