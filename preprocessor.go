package orgmode

import (
	"strings"
)

func Preprocess(content string) *OrgRoot {

	org := &OrgRoot{Settings: make(map[string]string)}

	// start of line
	lpos := 0
	line := ""
	for lookahead := range content {

		if content[lookahead] == LF {
			// indentation not important for document & option keywords
			line = content[lpos:lookahead]
			lpos = lookahead
		}

		if DefaultOrgRegexpMap["inbuffersetting"].MatchString(line) {
			m := DefaultOrgRegexpMap["inbuffersetting"].FindStringSubmatch(line)

			// TODO: Warn on overrides?
			org.Settings[strings.ToUpper(m[1])] = strings.TrimSpace(m[2])
		}
	}

	// TODO: Generate the new regexes for the Org document
	org.RegexpMap = DefaultOrgRegexpMap

	// TODO: Filter out affiliated keywords from settings
	return org
}
