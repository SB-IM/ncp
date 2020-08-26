package ncpio

import (
	"regexp"
)

type Rule struct {
	Regexp string `json:"regexp"`
	Invert bool   `json:"invert"`
}

func Filter(rules []Rule, data []byte) bool {
	for _, rule := range rules {
		matched, err := regexp.Match(rule.Regexp, data)
		if err != nil {
			continue
		}

		// Invert result
		if rule.Invert {
			matched = !matched
		}

		if matched {
			return true
		}
	}
	return false
}
