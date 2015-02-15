package orgmode

import (
	"regexp"
)

// Consider TODO DONE at the beginning only,
// later it would be overriden by the `#+TODO:` in buffer setting
const headlineRegexString = "^(\\*+)(?: +(CANCELED|DONE|TODO))?(?: +(\\[#.\\]))?(?: +(.*?))?(?:(:[a-zA-Z0-9_@#%:]+:))??[ \\t]*$"

// Compile static regexes, used during the first pass
var (
	DefaultOrgRegexpMap = map[string]*regexp.Regexp{
		"inbuffersetting":       regexp.MustCompile(`#\+(\w+): (.*)`),
		"headline":              regexp.MustCompile(`^(\*+)(?: +(CANCELED|DONE|TODO))?(?: +(\[#.\]))?(?: +(.*?))?(?:(:[a-zA-Z0-9_@#%:]+:))??[ \t]*$`),
		"beginKeywordOrComment": regexp.MustCompile(`^[ \t]*#`),
		"beginBlock":            regexp.MustCompile(`(?i)\+BEGIN_(CENTER|COMMENT|EXAMPLE|QUOTE|SRC|VERSE)`),
		"endBlock":              regexp.MustCompile(`(?i)\+END_(CENTER|COMMENT|EXAMPLE|QUOTE|SRC|VERSE)`),
		"blockInfo":             regexp.MustCompile(`[ \t]+([^ \f\t\n\r\v]+)[ \t]*([^\":\n]*\"[^\"\n*]*\"[^\":\n]*|[^\":\n]*)([^\n]*)`),
		"blockName":             regexp.MustCompile(`(?i)(^[ \t]*)#\+NAME:[ \t]*(.*)$`),
		"blockHeaders":          regexp.MustCompile(`(?i)(^[ \t]*)#\+HEADERS?:[ \t]*([^\n]*)$`),
		"beginSrc":              regexp.MustCompile(`(?i)[ \t]*#\+BEGIN_SRC`),
		"endSrc":                regexp.MustCompile(`(?i)[ \t]*#\+END_SRC`),
	}
)
