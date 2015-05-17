package gostruct

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

// helpers

func doc(t *testing.T, body string) *goquery.Document {
	root, err := html.Parse(strings.NewReader(body))
	require.Nil(t, err)
	return goquery.NewDocumentFromNode(root)
}

func baseDoc(t *testing.T) *goquery.Document {
	return doc(t, "<p>hello</p>")
}

func makeServer(t *testing.T, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, body)
		}))
}

// noop cases

func TestEmptyStruct(t *testing.T) {
	s := struct{}{}

	assert.Nil(t, Populate(&s, baseDoc(t)))
}

func TestNoSelectorsStruct(t *testing.T) {
	s := struct {
		Foo string
	}{}

	assert.Nil(t, Populate(&s, baseDoc(t)))
	assert.Equal(t, "", s.Foo)
}

func TestUnexportedFieldsStruct(t *testing.T) {
	s := struct {
		foo string
	}{}

	assert.Nil(t, Populate(&s, baseDoc(t)))
}

// wrong arguments

func TestNotAPointer(t *testing.T) {
	assert.NotNil(t, Populate(struct{}{}, baseDoc(t)))
	assert.NotNil(t, Populate(42, baseDoc(t)))
}

func TestNotAStructPointer(t *testing.T) {
	n := 2

	assert.NotNil(t, Populate(&n, baseDoc(t)))
}

// string field

func TestEmptyStringElement(t *testing.T) {
	s := struct {
		Title string `gostruct:"#title"`
	}{"something"}

	assert.Nil(t, Populate(&s, doc(t, `<h1 id="title"></h1>`)))
	assert.Equal(t, "", s.Title)
}

func TestStringElement(t *testing.T) {
	s := struct {
		Title string `gostruct:"#title"`
	}{}

	assert.Nil(t, Populate(&s, doc(t, `<h1 id="title">hello</h1>`)))
	assert.Equal(t, "hello", s.Title)
}

func TestMultipleStringElements(t *testing.T) {
	s := struct {
		Letters string `gostruct:".c"`
	}{}

	assert.Nil(t, Populate(&s, doc(t, `<p class="c">H</p>x<p class="c">i</p>`)))
	assert.Equal(t, "Hi", s.Letters)
}

// bool field

func TestEmptyBoolElement(t *testing.T) {
	s := struct {
		Title bool `gostruct:"#title"`
	}{true}

	d := doc(t, `<h1 id="title"></h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, false, s.Title)
}

func TestEmptyBoolSelection(t *testing.T) {
	s := struct {
		Title bool `gostruct:".foo"`
	}{}

	d := doc(t, `<h1 id="title">hi</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, false, s.Title)
}

func TestBoolElement(t *testing.T) {
	s := struct {
		Title bool `gostruct:"#title"`
	}{}

	d := doc(t, `<h1 id="title">hi</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, true, s.Title)
}

// int field

func TestEmptyIntElement(t *testing.T) {
	s := struct {
		Title int `gostruct:"#title"`
	}{42}

	d := doc(t, `<h1 id="title"></h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, 0, s.Title)
}

func TestEmptyIntSelection(t *testing.T) {
	s := struct {
		Title int `gostruct:".t"`
	}{42}

	d := doc(t, `<h1 id="title">37</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, 0, s.Title)
}

func TestIntWrongValue(t *testing.T) {
	s := struct {
		Title int `gostruct:"#title"`
	}{42}

	d := doc(t, `<h1 id="title">two</h1>`)

	assert.NotNil(t, Populate(&s, d))
	assert.Equal(t, 42, s.Title, "fields shouldn't be changed on parsing error")
}

func TestIntMultipleElements(t *testing.T) {
	s := struct {
		Count int `gostruct:"p"`
	}{}

	d := doc(t, `<p>42</p><p>17</p><p>1034</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, 42, s.Count, "int fields should only take the first element")
}

func TestNegativeInt(t *testing.T) {
	s := struct {
		Count int `gostruct:"p"`
	}{}

	d := doc(t, `<p>-42</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, -42, s.Count)
}

func TestInt(t *testing.T) {
	s := struct {
		Count int `gostruct:"p"`
	}{}

	d := doc(t, `<p>42</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, 42, s.Count)
}

func TestInt8(t *testing.T) {
	s := struct {
		Count int8 `gostruct:"p"`
	}{}

	d := doc(t, `<p>42</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, int8(42), s.Count)
}

// uint field

func TestEmptyUintElement(t *testing.T) {
	s := struct {
		Count uint `gostruct:"p"`
	}{}

	d := doc(t, `<p></p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, uint64(0), s.Count)
}

func TestUintNegativeInt(t *testing.T) {
	s := struct {
		Count uint `gostruct:"p"`
	}{}

	d := doc(t, `<p>-42</p>`)

	assert.NotNil(t, Populate(&s, d))
}

func TestUint(t *testing.T) {
	s := struct {
		Count uint `gostruct:"p"`
	}{}

	d := doc(t, `<p>42</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, uint(42), s.Count)
}

// float field

func TestFloatFirstElement(t *testing.T) {
	s := struct {
		Count float64 `gostruct:"p"`
	}{}

	d := doc(t, `<p>42.6</p><p>25.7</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, float64(42.6), s.Count)
}

func TestEmptyFloatElement(t *testing.T) {
	s := struct {
		Count float64 `gostruct:"p"`
	}{}

	d := doc(t, `<p></p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, float64(0), s.Count)
}

func TestFloat(t *testing.T) {
	s := struct {
		Count float64 `gostruct:"p"`
	}{}

	d := doc(t, `<p>42.6</p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, float64(42.6), s.Count)
}

// sub-structs

func TestEmptySubStruct(t *testing.T) {
	s := struct {
		X struct{} `gostruct:"#foo"`
	}{}

	assert.Nil(t, Populate(&s, doc(t, `<p id="#foo">x</p>`)))
}

func TestUnexportedSubStruct(t *testing.T) {
	s := struct {
		x struct {
			Foo string `gostruct:"#foo"`
		} `gostruct:".x"`
	}{}

	d := doc(t, `<p class="x"><span id="foo">Bar</span></p>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, "", s.x.Foo)
}

func TestSubStruct(t *testing.T) {
	s := struct {
		Header struct {
			Title string `gostruct:"h1"`
		} `gostruct:"header"`
	}{}

	d := doc(t, `<header><h1>Good</h1></header><h1>Bad</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, "Good", s.Header.Title)
}

func TestSubStructPtr(t *testing.T) {
	s := struct {
		Header *struct {
			Title string `gostruct:"h1"`
		} `gostruct:"header"`
	}{}

	d := doc(t, `<header><h1>Good</h1></header><h1>Bad</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, "Good", s.Header.Title)
}

// fetching test

func TestFetchWrongURL(t *testing.T) {
	s := struct{}{}

	assert.NotNil(t, Fetch(&s, "::not an URL"))
}

func TestFetch(t *testing.T) {
	server := makeServer(t, `<h1>Hello</h1>`)
	defer server.Close()

	s := struct {
		Title string `gostruct:"h1"`
	}{}

	assert.Nil(t, Fetch(&s, server.URL))
	assert.Equal(t, "Hello", s.Title)
}

// byte slice tests

func TestPopulateByteSliceEmptyString(t *testing.T) {
	s := struct {
		Title []byte `gostruct:"h1"`
	}{}

	d := doc(t, `<h1></h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, []byte(""), s.Title)
}

func TestPopulateByteSlice(t *testing.T) {
	s := struct {
		Title []byte `gostruct:"h1"`
	}{}

	d := doc(t, `<h1>Hello</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, []byte("Hello"), s.Title)
}

// slice tests

func TestPopulateSliceEmptySelection(t *testing.T) {
	s := struct {
		Names []string `gostruct:"li"`
	}{}

	d := doc(t, `<ol></ol>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, []string{}, s.Names)
}

func TestPopulateStringSlice(t *testing.T) {
	s := struct {
		Names []string `gostruct:"li"`
	}{}

	d := doc(t, `<ol><li>A</li><li>B</li><li>C</li></ol>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, []string{"A", "B", "C"}, s.Names)
}

// general tests

func TestPopulatePointer(t *testing.T) {
	s := &struct {
		Title string `gostruct:"h1"`
	}{}

	d := doc(t, `<h1>This is a test</h1>`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, "This is a test", s.Title)
}

func TestMultipleFields(t *testing.T) {
	s := struct {
		Title  string        `gostruct:"h1"`
		Desc   string        `gostruct:"p"`
		PubAge time.Duration `gostruct:".pub-age"`
	}{}

	d := doc(t, `
		<h1>This is a test</h1>
		<p>This is its description</p>
		<div class="metadata">
			Author: <a href="/foo">Foo</a><br/>
			Published <span class="pub-age">1h30m</span> ago.
		</div>
	`)

	assert.Nil(t, Populate(&s, d))
	assert.Equal(t, "This is a test", s.Title)
	assert.Equal(t, "This is its description", s.Desc)
	assert.Equal(t, "1h30m0s", s.PubAge.String())
}
