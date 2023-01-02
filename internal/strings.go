package internal

import (
	"regexp"
	"strconv"
)

func ParseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		HandleMainError(err)
	}
	return n
}

func NamedCaptureGroups(regexp *regexp.Regexp, s string) map[string]string {
	match := regexp.FindStringSubmatch(s)
	len := len(match)

	ret := make(map[string]string, len-1)
	for i, name := range regexp.SubexpNames() {
		if name != "" {
			ret[name] = match[i]
		}
	}

	return ret
}
