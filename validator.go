// Package validator
package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type validator struct {
}

func validate(value interface{}, fieldName string, tags []string, errBag url.Values) error {

	isRequired := false

	for i := 0; i < len(tags); i++ {

		rl := strings.Split(tags[i], ":")

		fn := rl[0]

		if fn == "required" {
			isRequired = true
		}

		f, ok := rules[fn]
		if !ok {
			continue
		}

		err := f(value, fieldName, tags[i], isRequired)
		if err == nil {
			continue
		}

		errBag.Add(fieldName, err.Error())

	}

	return nil

}

func validateStruct(v reflect.Value, parentField string) url.Values {
	errBag := url.Values{}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := t.Field(i)

		fn := fi.Tag.Get("json")
		fv := fi.Tag.Get("valid")

		if fv == "" || fv == "-" {
			continue
		}

		tags := strings.Split(fv, "|")

		if parentField != "" {
			fn = fmt.Sprintf("%s.%s", parentField, fn)
		}

		validate(v.Field(i).Interface(), fn, tags, errBag)

		switch v.Field(i).Kind() {
		case reflect.Struct:
			eb := validateStruct(v.Field(i), fn)
			mergeKeys(errBag, eb)
		case reflect.Slice:
		case reflect.Ptr:
			ptrRef := reflect.Indirect(v.Field(i))
			switch ptrRef.Kind() {
			case reflect.Struct:
				eb := validateStruct(ptrRef, fn)
				mergeKeys(errBag, eb)
			}
		}

	}

	return errBag
}

func New() *validator {
	return &validator{}
}

func (vl *validator) ValidateStruct(input interface{}) url.Values {
	var errBag url.Values

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		errBag.Set(`_error`, `validator: invalid input type`)
		return errBag
	}


	return validateStruct(val, "")
}
