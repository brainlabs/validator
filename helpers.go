// Package validator
package validator

import (
	"fmt"
	"strings"
)

// containsRequiredField check rules contain any required field
func isContainRequiredField(rules []string) bool {
	for _, rule := range rules {
		if rule == "required" {
			return true
		}
	}
	return false
}

// isRuleExist check if the provided rule name is exist or not
func isRuleExist(rule string) bool {
	if strings.Contains(rule, ":") {
		rule = strings.Split(rule, ":")[0]
	}
	extendedRules := []string{"size", "mime", "ext"}
	for _, r := range extendedRules {
		if r == rule {
			return true
		}
	}
	if _, ok := rulesMap[rule]; ok {
		return true
	}

	return false
}

// toString force data to be string
func toString(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		str = fmt.Sprintf("%#v", v)
	}
	return str
}
