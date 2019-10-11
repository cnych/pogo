package strs

import "regexp"

func MatchRegexp(pattern, content string, index int) string {
	compile := regexp.MustCompile(pattern)
	matches := compile.FindStringSubmatch(content)
	for i, v := range matches {
		if i == index {
			return v
		}
	}
	return ""
}


func MatchRegexpOf1(pattern, content string) string {
	return MatchRegexp(pattern, content, 1)
}

