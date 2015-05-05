package gostruct

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

func Fetch(target interface{}, url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	return Populate(target, doc)
}

func Populate(target interface{}, doc *goquery.Document) error {
	value := reflect.ValueOf(target)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value '%s' is not a pointer", target)
	}

	elem := value.Elem()

	// from https://github.com/vrischmann/envconfig/blob/master/envconfig.go#L38
	switch elem.Kind() {
	case reflect.Ptr:
		elem.Set(reflect.New(elem.Type().Elem()))
		return populateStruct(elem.Elem(), doc.Selection)
	case reflect.Struct:
		return populateStruct(elem, doc.Selection)
	default:
		return fmt.Errorf("value '%s' must be a pointer on struct", target)
	}
}

func populateStruct(target reflect.Value, doc *goquery.Selection) (err error) {
	fieldsCount := target.NumField()
	targetType := target.Type()

	for i := 0; i < fieldsCount; i++ {
		field := target.Field(i)
		sel := targetType.Field(i).Tag.Get("gostruct")
		if sel == "" {
			continue
		}

		subdoc := doc.Find(sel)

		// from https://github.com/vrischmann/envconfig/blob/master/envconfig.go#L87
	doPopulate:
		switch field.Kind() {
		case reflect.Ptr:
			field.Set(reflect.New(field.Type().Elem()))
			field = field.Elem()
			goto doPopulate
		case reflect.Struct:
			err = populateStruct(field, subdoc)
		default:
			err = setField(field, subdoc)
		}

		if err != nil {
			break
		}
	}

	return
}

func setField(field reflect.Value, doc *goquery.Selection) error {
	if doc.Length() == 0 {
		return nil
	}

	return errors.New("TODO")
}
