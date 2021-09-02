// Package validator
package validator

import (
	"fmt"
	"reflect"
	"strings"
)

type rule func(value interface{}, fieldName, tagRule string, isRequired bool) error

var (
	rules = map[string]rule{
		"required":    Required,
		"numeric":     ValidNumeric,
		"float":       ValidFloat,
		"max":         ValidMax,
		"min":         ValidMin,
		"alpha_num":   ValidAlphaNum,
		"alpha_space": ValidAlphaSpace,
		"alpha_dash":  ValidAlphaDash,
		"email":       ValidEmail,
		"uuid":        ValidUUID,
		"uuid3":       ValidUUID3,
		"uuid4":       ValidUUID4,
		"uuid5":       ValidUUID5,
		"url":         ValidURL,
		"credit_card": ValidCreditCard,
		"latitude":    ValidLatitude,
		"longitude":   ValidLongitude,
		"mac_address": ValidMacAddress,
		"coordinate":  ValidCoordinate,
		"ip":          ValidIP,
		"ipv4":        ValidIPV4,
		"ipv6":        ValidIPV6,
		"imei":        ValidIMEI,
		"hex_color":   ValidHexColor,
		"isbn10":      ValidISBN10,
		"isbn13":      ValidISBN13,
		"json":        ValidJSON,
		"bool":        ValidBoolean,
		"in":          ValidIn,
		"id_phone":    ValidIndonesianPhoneNumber,
	}
)

var rulesMap = map[string]interface{}{}

var funcSignature = []string{
	"func(interface {}, string, string) error",
	"func(interface {}, string, string, string) error",
	"func(interface {}, string, interface {}, string, string, string) error",
	"func(interface {}, string, string, string, string) error",
}

func AddNewRule(name string, fn func(field string, rule string, message string, value interface{}) error) error {
	if isRuleExist(name) {
		return fmt.Errorf("validator: %s is already defined in rules %+v", name)
	}

	rulesMap[name] = fn
	return nil
}

func AddNeFunc(key string, f interface{}) error {
	fValue := reflect.ValueOf(f)
	if fValue.Kind() != reflect.Func {
		return fmt.Errorf("please provide a function typed argument")
	}

	notFound := true
	for _, s := range funcSignature {
		if fValue.Type().String() == s {
			notFound = false
			break
		}
	}

	if notFound {
		return fmt.Errorf("function signature is not accepted")
	}

	_, ok := rulesMap[key]

	if ok {
		return fmt.Errorf("function already registered")
	}

	return nil
}

func Required(v interface{}, key, rule string, isRequired bool) error {
	msg := `The %s field is required`

	if isEmpty(v) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidEmail(v interface{}, key, rule string, isRequired bool) error {
	msg := `The %s field should be a valid email address`

	if isEmpty(v) && !isRequired {
		return nil
	}

	vs := ToString(v)

	if isEmail(vs) {
		return nil
	}

	return fmt.Errorf(msg, key)
}

func ValidNumeric(v interface{}, key, rule string, isRequired bool) error {
	msg := `The %s field should be a valid numeric`

	if isEmpty(v) && !isRequired {
		return nil
	}

	vs := ToString(v)

	if isNumeric(vs) {
		return nil
	}

	return fmt.Errorf(msg, key)
}

func ValidMin(v interface{}, key, rule string, isRequired bool) error {

	msgInt := `The %s field should be greater than or equal %s`
	msgStr := `The %s field should be minimum length %s`

	if isEmpty(v) && !isRequired {
		return nil
	}

	min := strings.TrimPrefix(rule, "min:")

	if !isNumeric(min) {
		return fmt.Errorf(`The %s field invalid rule format %s`, key, rule)
	}

	vs := ToString(v)

	cm, _ := ToInt64(min)

	if isNumeric(vs) {
		vi, _ := ToInt64(vs)

		if vi < cm {
			return fmt.Errorf(msgInt, key, min)
		}

		return nil
	}

	if len(vs) >= int(cm) {
		return nil
	}

	return fmt.Errorf(msgStr, key, min)
}

func ValidMax(v interface{}, key, rule string, isRequired bool) error {

	msgInt := `The %s field should be less than or equal %s`
	msgStr := `The %s field should be maximum length %s`

	if isEmpty(v) && !isRequired {
		return nil
	}

	max := strings.TrimPrefix(rule, "max:")

	if !isNumeric(max) {
		return fmt.Errorf(`The %s field invalid rule format %s`, key, rule)
	}

	vs := ToString(v)

	cm, _ := ToInt64(max)

	if isNumeric(vs) {
		vi, _ := ToInt64(vs)
		if vi > cm || (isRequired && isEmpty(v)) {
			return fmt.Errorf(msgInt, key, max)
		}

		return nil
	}

	if len(vs) > int(cm) || (isRequired && isEmpty(v)) {
		return fmt.Errorf(msgStr, key, max)
	}

	return nil
}

func ValidAlphaNum(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should contain: [a-zA-Z0-9]`

	s := ToString(v)

	if !isAlphaNumeric(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidAlphaDash(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should contain: [a-zA-Z0-9], underscore (_), dash (-)`

	s := ToString(v)

	if !isAlphaDash(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidAlphaSpace(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should contain: [a-zA-Z0-9], underscore (_), space`

	s := ToString(v)

	if !isAlphaSpace(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidAlpha(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should contain: [a-zA-Z]`

	s := ToString(v)

	if !isAlpha(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidURL(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid URL`

	s := ToString(v)

	if !isURL(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidMacAddress(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid mac address`

	s := ToString(v)

	if !isMacAddress(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidFloat(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be float number`

	s := ToString(v)

	if !isFloat(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidUUID(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be uuid`

	s := ToString(v)

	if !isUUID(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidUUID3(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be uuid3`

	s := ToString(v)

	if !isUUID3(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidUUID4(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be uuid4`

	s := ToString(v)

	if !isUUID4(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidUUID5(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be uuid5`

	s := ToString(v)

	if !isUUID5(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidCreditCard(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be credit card number`

	s := ToString(v)

	if !isCreditCard(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidLatitude(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid latitude`

	s := ToString(v)

	if !isLatitude(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidLongitude(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid longitude`

	s := ToString(v)

	if !isLongitude(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidCoordinate(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid coordinate`

	s := ToString(v)

	if !isCoordinate(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidIP(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid IP Address`

	s := ToString(v)

	if !isIP(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidIPV4(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid IPV4 Address`

	s := ToString(v)

	if !isIPV4(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidIPV6(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid IPV6 Address`

	s := ToString(v)

	if !isIPV6(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidIMEI(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid IMEI`

	s := ToString(v)

	if !isIMEI(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidHexColor(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid hexa color`

	s := ToString(v)

	if !isHexColor(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidISBN10(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid ISBN10`

	s := ToString(v)

	if !isISBN10(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidISBN13(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid ISBN13`

	s := ToString(v)

	if !isISBN13(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidJSON(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid JSON`

	s := ToString(v)

	if !isJSON(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidBoolean(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should be valid boolean type`

	s := ToString(v)

	if !isBoolean(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}

func ValidIn(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	contain := strings.Trim(rule, "in:")

	haystack := strings.Split(contain, ",")

	msg := `The %s field should be contain in: %s`

	s := ToString(v)

	if !isIn(haystack, s) {
		return fmt.Errorf(msg, key, contain)
	}

	return nil
}

func ValidIndonesianPhoneNumber(v interface{}, key, rule string, isRequired bool) error {
	if isEmpty(v) && !isRequired {
		return nil
	}

	msg := `The %s field should valid mobile phone number`

	s := ToString(v)

	if !isIndonesiaPhoneNumber(s) {
		return fmt.Errorf(msg, key)
	}

	return nil
}
