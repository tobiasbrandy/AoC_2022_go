package regext

import "regexp"

func NamedCaptureGroups(regexp *regexp.Regexp, s string) map[string]string {
	match := regexp.FindStringSubmatch(s)
	l := len(match)
	if l == 0 {
		return make(map[string]string, 0)
	}

	ret := make(map[string]string, l-1)
	for i, name := range regexp.SubexpNames() {
		if name != "" {
			ret[name] = match[i]
		}
	}

	return ret
}
