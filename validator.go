// Package validator
package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	defaultTagRule  = `valid`
	defaultTagField = `json`
)

type Option func(*Validator)

type Validator struct {
	TagField string
	TagRule  string
}

// OptionTagField option tag field
func OptionTagField(tag string) Option {
	return func(v *Validator) {
		v.TagField = tag
	}
}

// OptionTagValidationRule option tag validation rule
func OptionTagRule(tag string) Option {
	return func(v *Validator) {
		v.TagRule = tag
	}
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

func validateStruct(v reflect.Value, parentField, tagField, tagRule string) url.Values {
	errBag := url.Values{}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := t.Field(i)

		tf := fi.Tag.Get(tagField)
		tr := fi.Tag.Get(tagRule)

		if tr == "" || tr == "-" {
			continue
		}

		tags := strings.Split(tr, "|")

		if parentField != "" {
			tf = fmt.Sprintf("%s.%s", parentField, tf)
		}

		validate(v.Field(i).Interface(), tf, tags, errBag)

		switch v.Field(i).Kind() {
		case reflect.Struct:
			eb := validateStruct(v.Field(i), tf, tagField, tagRule)
			mergeKeys(errBag, eb)
		case reflect.Slice:
		case reflect.Ptr:
			ptrRef := reflect.Indirect(v.Field(i))
			switch ptrRef.Kind() {
			case reflect.Struct:
				eb := validateStruct(ptrRef, tf, tagField, tagRule)
				mergeKeys(errBag, eb)
			}
		}

	}

	return errBag
}

func New(options ...Option) *Validator {
	x := &Validator{}

	for _, opt := range options {
		opt(x)
	}

	if x.TagRule == "" {
		x.TagRule = defaultTagRule
	}

	if x.TagField == "" {
		x.TagField = defaultTagField
	}

	return x
}

func (vl *Validator) ValidateStruct(input interface{}) url.Values {
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

	return validateStruct(val, "", vl.TagField, vl.TagRule)
}
