package gostruct

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func doc(t *testing.T, body string) *goquery.Document {
	root, err := html.Parse(strings.NewReader(body))
	require.Nil(t, err)
	return goquery.NewDocumentFromNode(root)
}

func baseDoc(t *testing.T) *goquery.Document {
	return doc(t, "<p>hello</p>")
}

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

func TestNotAPointer(t *testing.T) {
	assert.NotNil(t, Populate(struct{}{}, baseDoc(t)))
	assert.NotNil(t, Populate(42, baseDoc(t)))
}
