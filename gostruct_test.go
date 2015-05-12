package gostruct

import (
	"strings"
	"testing"

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
