package gostruct

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

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
		default:
			err = setField(field, subdoc)
		}

		if err != nil {
			break
		}
	}

	return
}

var (
	durationType  = reflect.TypeOf(new(time.Duration)).Elem()
	byteSliceType = reflect.TypeOf([]byte(nil))
)

func isDurationField(t reflect.Type) bool {
	return t.AssignableTo(durationType)
}

func setField(field reflect.Value, doc *goquery.Selection) error {
	if !field.CanSet() {
		// unexported field: don't do anything
		return nil
	}

	ftype := field.Type()
	kind := ftype.Kind()

	// types which take the whole selection
	switch kind {
	case reflect.Struct:
		return populateStruct(field, doc)
	case reflect.Slice:
		if ftype == byteSliceType {
			return setByteSliceValue(field, doc)
		}
		return setSliceValue(field, doc)
	case reflect.String:
		return setStringValue(field, doc)
	case reflect.Bool:
		return setBoolValue(field, doc)
	}

	text := doc.First().Text()

	// types which take only the first element's text

	if isDurationField(ftype) {
		return setDurationValue(field, text)
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntValue(field, text)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintValue(field, text)
	case reflect.Float32, reflect.Float64:
		return setFloatValue(field, text)
	default:
		return fmt.Errorf("Unsupported field type: '%v'", ftype)
	}
}

func setStringValue(field reflect.Value, sel *goquery.Selection) error {
	field.SetString(sel.Text())
	return nil
}

func setBoolValue(field reflect.Value, sel *goquery.Selection) error {
	// this one is tricky because there are multiple possible interpretations:
	// - set to true only if there are elements matching the selector
	// - set to true if the selection's text is not empty (this is what we're
	//   doing here)
	// - set to the resulting value of `strconf.ParseBool` called on the
	//   selection's text
	field.SetBool(sel.Text() != "")
	return nil
}

func setIntValue(field reflect.Value, s string) error {
	if s == "" {
		field.SetInt(0)
		return nil
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		field.SetInt(val)
	}

	return err
}

func setUintValue(field reflect.Value, s string) error {
	if s == "" {
		field.SetUint(0)
		return nil
	}

	val, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		field.SetUint(val)
	}

	return err
}

func setFloatValue(field reflect.Value, s string) error {
	if s == "" {
		field.SetFloat(0)
		return nil
	}

	val, err := strconv.ParseFloat(s, 64)
	if err == nil {
		field.SetFloat(val)
	}

	return err
}

func setDurationValue(field reflect.Value, s string) error {
	val, err := time.ParseDuration(s)
	if err == nil {
		field.SetInt(int64(val))
	}

	return err
}

// this one is like setStringValue except that we convert the string in a byte
// slice
func setByteSliceValue(field reflect.Value, sel *goquery.Selection) error {
	field.SetBytes([]byte(sel.Text()))
	return nil
}

func setSliceValue(field reflect.Value, sel *goquery.Selection) error {
	count := sel.Length()

	eltype := field.Type().Elem()
	capacity := field.Cap()

	if count > capacity {
		capacity = count
	}

	slice := reflect.MakeSlice(field.Type(), 0, capacity)

	var err error

	sel.EachWithBreak(func(i int, subSel *goquery.Selection) bool {
		el := reflect.New(eltype).Elem()

		if err = setField(el, subSel); err != nil {
			return false
		}

		slice = reflect.Append(slice, el)

		return true
	})

	if err != nil {
		return err
	}

	field.Set(slice)

	return nil
}
