package str

import (
	"strings"
	"unicode"
)

func Slug(s string) string {
	_s := ""
	for _, r := range strings.ReplaceAll(strings.TrimSpace(s), "	", " ") {
		if unicode.IsLetter(r) {
			_s = _s + string(r)
			continue
		}
		if unicode.IsDigit(r) {
			_s = _s + string(r)
			continue
		}
		_s = strings.TrimSuffix(_s, "-")
		_s = _s + "-"
	}
	return strings.ToLower(strings.TrimSuffix(_s, "-"))
}
