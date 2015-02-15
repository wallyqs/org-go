package orgexport

import (
	"fmt"
	. "org-mode"
	"strconv"
	"strings"
)

func ToSexp(o interface{}, indentation ...int) string {

	sexp := ""

	indentCounter := 0
	if len(indentation) > 0 {
		indentCounter = indentation[0]
	}
	indent := strings.Repeat(" ", indentCounter)

	// Should be pretty printing here the object
	switch oo := o.(type) {

	case *OrgRoot:

		sexp += fmt.Sprintf("%s (section (:elements %v)", indent, len(oo.Elements))
		indentCounter += 2
		indent = strings.Repeat(" ", indentCounter)
		for k, v := range oo.Settings {
			sexp += fmt.Sprintf("\n%s(keyword (:key %s :value %s ))", indent, strconv.Quote(k), strconv.Quote(v))
		}

		for _, e := range oo.Elements {
			sexp += ToSexp(e, indentCounter)
		}

	case *OrgHeadline:

		sexp += fmt.Sprintf("\n%s(headline (:raw-value %s :level %v :elements %v :state %s :priority %s :tags %s :title %s)", indent, strconv.Quote(oo.RawContent), oo.Level, len(oo.Elements), strconv.Quote(oo.State), strconv.Quote(oo.Priority), strconv.Quote(oo.Tags), strconv.Quote(oo.Title))
		indentCounter += 2
		indent = strings.Repeat(" ", indentCounter)

		for _, e := range oo.Elements {
			sexp += ToSexp(e, indentCounter)
		}

	case *OrgSrcBlock:

		var headerArgs string
		for key, value := range oo.Headers {
			headerArgs = fmt.Sprintf("%s %s %s", headerArgs, key, value)
		}

		sexp += fmt.Sprintf("\n%s(src-block (:language %s :switches %s :parameters %s :value %s)", indent, strconv.Quote(oo.Lang), strconv.Quote(oo.Switches), strconv.Quote(headerArgs), strconv.Quote(oo.RawContent))

	case *OrgParagraph:

		sexp += fmt.Sprintf("\n%s(paragraph (:raw-value %s)", indent, strconv.Quote(oo.RawContent))

	}

	sexp += ")"

	return sexp
}
