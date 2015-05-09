package gostruct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// some parts of this code are stolen^Winspired from
//   https://github.com/vrischmann/envconfig

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

func setField(field reflect.Value, doc *goquery.Selection) (err error) {
	if !field.CanSet() {
		// unexported field: don't do anything
		return nil
	}

	text := doc.Text()

	switch field.Type().Kind() {
	case reflect.Bool:
		err = setBoolValue(field, text)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = setIntValue(field, text)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = setUintValue(field, text)
	case reflect.Float32, reflect.Float64:
		err = setFloatValue(field, text)
	case reflect.String:
		err = setStringValue(field, text)
	default:
		errors.New("TODO")
	}

	return
}

func setBoolValue(field reflect.Value, s string) error {
	// this one is tricky because there are multiple possible interpretations:
	// - set to true only if there are elements matching the selector
	// - set to true if the selection's text is not empty (this is what we're
	//   doing here)
	// - set to the resulting value of `strconf.ParseBool` called on the
	//   selection's text
	field.SetBool(s != "")
	return nil
}

func setIntValue(field reflect.Value, s string) error {
	val, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		field.SetInt(val)
	}

	return err
}

func setUintValue(field reflect.Value, s string) error {
	val, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		field.SetUint(val)
	}

	return err
}

func setFloatValue(field reflect.Value, s string) error {
	val, err := strconv.ParseFloat(s, 64)
	if err == nil {
		field.SetFloat(val)
	}

	return err
}

func setStringValue(field reflect.Value, s string) error {
	field.SetString(s)
	return nil
}
