// Package validator
package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// funcval, errorMessage, key
// funcval, format, errorMessage, key
// funcval, errMessage, key1, key2

type dataTag struct {
	funcVal        string
	errorMessage   string
	format         string
	compareKey     string
	compareValue   string
	dateLayout     string
	acceptedValues string
}

// fetchDataTag idx must always starts from -1
func fetchDataTag(input string, idx int, dataTags []*dataTag) []*dataTag {
	if input == "" {
		return dataTags
	}
	tagsSplits := strings.Split(input, "|")
	if len(tagsSplits) > 1 {
		for i := 0; i < len(tagsSplits); i++ {
			dataTags = append(dataTags, &dataTag{})
		}

		for i := 0; i < len(tagsSplits); i++ {
			fetchDataTag(tagsSplits[i], i, dataTags)
		}

		return dataTags
	}

	aTagSplits := strings.Split(input, ",")
	if len(aTagSplits) > 1 {
		if len(dataTags) == 0 {
			dataTags = append(dataTags, &dataTag{})
			idx = 0
		}
		for i := 0; i < len(aTagSplits); i++ {
			fetchDataTag(aTagSplits[i], idx, dataTags)
		}
		return dataTags
	}

	splits := strings.Split(input, ":")
	if len(dataTags) == 0 {
		dataTags = append(dataTags, &dataTag{})
		idx = 0
	}

	if len(splits) > 1 {
		iTag := dataTags[idx]
		if iTag == nil {
			iTag = &dataTag{}
			dataTags[idx] = iTag
		}
		switch splits[0] {
		case "funcVal":
			iTag.funcVal = splits[1]
		case "errorMessage":
			iTag.errorMessage = splits[1]
		case "format":
			iTag.format = splits[1]
		case "compareValue":
			iTag.compareValue = splits[1]
		case "compareKey":
			iTag.compareKey = splits[1]
		case "dateLayout":
			iTag.dateLayout = splits[1]
		case "values":
			iTag.acceptedValues = splits[1]
		}
		fetchDataTag("", idx, dataTags)
	}

	return dataTags

}

func DumpToString(v interface{}) string {

	str, ok := v.(string)
	if !ok {
		buff := &bytes.Buffer{}
		json.NewEncoder(buff).Encode(v)
		return buff.String()
	}

	return str
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}

	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	return v.Interface()
}

func mergeKeys(left, right url.Values) url.Values {

	if len(right) == 0 {
		return left
	}
	for key, rightVal := range right {
		if _, present := left[key]; present {
			//then we don't want to replace it - recurse
			left[key] = append(left[key], rightVal...)
			continue
		}
		// key not in left so we can just shove it in
		left[key] = rightVal
	}

	return left
}

// Match regular expression validation
func Match(value interface{}, key, format, msg string) error {

	rgx, e := regexp.Compile(format)

	if e != nil {
		return fmt.Errorf("%s invalid rule regular expression %s: %s", key, format, e.Error())
	}

	val, ok := value.(string)

	if !ok {
		return fmt.Errorf("invalid type, expected string found %s", reflect.TypeOf(value))
	}

	if !rgx.MatchString(val) {
		if msg == "" {
			return fmt.Errorf("%s has invalid format value", key)
		}

		return fmt.Errorf(msg)
	}

	return nil
}

// isEmpty check a type is Zero
func isEmpty(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(v.String())) == 0
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

// ToString converts a value to string.
func ToString(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return strconv.FormatInt(int64(value), 10)
	case int8:
		return strconv.FormatInt(int64(value), 10)
	case int16:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(int64(value), 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'g', -1, 64)
	case float64:
		return strconv.FormatFloat(float64(value), 'g', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	default:
		return fmt.Sprintf("%+v", value)
	}
}

// ToInt64 cast interface to an int64 type.
func ToInt64(i interface{}) (int64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int:
		return int64(s), nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

// isAlpha check the input is letters (a-z,A-Z) or not
func isAlpha(str string) bool {
	return regexAlpha.MatchString(str)
}

// isAlphaDash check the input is letters, number with dash and underscore
func isAlphaDash(str string) bool {
	return regexAlphaDash.MatchString(str)
}

// isAlphaSpace check the input is letters, number with dash and underscore
func isAlphaSpace(str string) bool {
	return regexAlphaSpace.MatchString(str)
}

// isAlphaNumeric check the input is alpha numeric or not
func isAlphaNumeric(str string) bool {
	return regexAlphaNumeric.MatchString(str)
}

// isBoolean check the input contains boolean type values
// in this case: "0", "1", "true", "false", "True", "False"
func isBoolean(str string) bool {
	bools := []string{"0", "1", "true", "false", "True", "False"}
	for _, b := range bools {
		if b == str {
			return true
		}
	}
	return false
}

//isCreditCard check the provided card number is a valid
//  Visa, MasterCard, American Express, Diners Club, Discover or JCB card
func isCreditCard(card string) bool {
	return regexCreditCard.MatchString(card)
}

// isCoordinate is a valid Coordinate or not
func isCoordinate(str string) bool {
	return regexCoordinate.MatchString(str)
}

// isCSSColor is a valid CSS color value (hex, rgb, rgba, hsl, hsla) etc like #909, #00aaff, rgb(255,122,122)
func isCSSColor(str string) bool {
	return regexCSSColor.MatchString(str)
}

// isDate check the date string is valid or not
func isDate(date string) bool {
	return regexDate.MatchString(date)
}

// isDateDDMMYY check the date string is valid or not
func isDateDDMMYY(date string) bool {
	return regexDateDDMMYY.MatchString(date)
}

// isEmail check a email is valid or not
func isEmail(email string) bool {
	return regexEmail.MatchString(email)
}

// isFloat check the input string is a float or not
func isFloat(str string) bool {
	return regexFloat.MatchString(str)
}

// isIn check if the niddle exist in the haystack
func isIn(haystack []string, niddle string) bool {
	for _, h := range haystack {
		if h == niddle {
			return true
		}
	}
	return false
}

// isJSON check wheather the input string is a valid json or not
func isJSON(str string) bool {
	var data interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return false
	}
	return true
}

// isNumeric check the provided input string is numeric or not
func isNumeric(str string) bool {
	return regexNumeric.MatchString(str)
}

// isMacAddres check the provided string is valid Mac Address or not
func isMacAddress(str string) bool {
	return regexMacAddress.MatchString(str)
}

// isLatitude check the provided input string is a valid latitude or not
func isLatitude(str string) bool {
	return regexLatitude.MatchString(str)
}

// isLongitude check the provided input string is a valid longitude or not
func isLongitude(str string) bool {
	return regexLongitude.MatchString(str)
}

// isIP check the provided input string is a valid IP address or not
func isIP(str string) bool {
	return regexIP.MatchString(str)
}

// isIPV4 check the provided input string is a valid IP address version 4 or not
// Ref: https://en.wikipedia.org/wiki/IPv4
func isIPV4(str string) bool {
	return regexIPV4.MatchString(str)
}

// isIPV6 check the provided input string is a valid IP address version 6 or not
// Ref: https://en.wikipedia.org/wiki/IPv6
func isIPV6(str string) bool {
	return regexIPV6.MatchString(str)
}

// isMatchedRegex match the regular expression string provided in first argument
// with second argument which is also a string
func isMatchedRegex(rxStr, str string) bool {
	rx := regexp.MustCompile(rxStr)
	return rx.MatchString(str)
}

// isURL check a URL is valid or not
func isURL(url string) bool {
	return regexURL.MatchString(url)
}

// isUUID check the provided string is valid UUID or not
func isUUID(str string) bool {
	return regexUUID.MatchString(str)
}

// isUUID3 check the provided string is valid UUID version 3 or not
func isUUID3(str string) bool {
	return regexUUID3.MatchString(str)
}

// isUUID4 check the provided string is valid UUID version 4 or not
func isUUID4(str string) bool {
	return regexUUID4.MatchString(str)
}

// isUUID5 check the provided string is valid UUID version 5 or not
func isUUID5(str string) bool {
	return regexUUID5.MatchString(str)
}

// isIMEI check the provided string is valid IMEI or not
func isIMEI(str string) bool {
	return regexIMEI.MatchString(str)
}

// isHexColor check the provided string is valid hexa color or not
func isHexColor(str string) bool {
	return regexHexColor.MatchString(str)
}

// isISBN10 check the provided string is valid ISBN10 or not
func isISBN10(str string) bool {
	return regexISBN10.MatchString(str)
}

// isISBN13 check the provided string is valid ISBN13 or not
func isISBN13(str string) bool {
	return regexISBN13.MatchString(str)
}


// isIndonesiaPhoneNumber check the provided string is valid indonesian phone number or not
func isIndonesiaPhoneNumber(str string) bool {
	return regexPhoneNumberID.MatchString(str)
}
