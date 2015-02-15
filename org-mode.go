package orgmode

import (
	"regexp"
)

const VERSION = "0.0.1"

const (
	ORG_HEADER_ARGUMENT_BEGIN = iota
	ORG_HEADER_ARGUMENT_KEYWORD_BEGIN
	ORG_HEADER_ARGUMENT_KEYWORD_END
	ORG_HEADER_ARGUMENT_VALUE_BEGIN
	ORG_HEADER_ARGUMENT_VALUE_END
)

const (
	TAB   = 9
	COLON = 58
	SPACE = 32
	LF    = 10
	CR    = 13
	NUL   = 0
)

type OrgElement struct {
	RawContent string
	Parent     interface{}
}

type OrgGreaterElement struct {
	Elements []interface{}
	OrgElement
}

type OrgRoot struct {
	OrgGreaterElement
	Settings  map[string]string
	Include   string
	RegexpMap map[string]*regexp.Regexp
	Level     int
}

type OrgHeadline struct {
	OrgGreaterElement
	Level    int
	State    string
	Priority string
	Title    string
	Tags     string // []string might be better
}

type OrgParagraph struct {
	OrgElement
}

type OrgEmptyLine struct {
	OrgElement
}

type OrgPropertiesDrawer struct {
	Options map[string]string
}

type OrgBlock struct {
	OrgElement
	Name string
}

type OrgSrcBlock struct {
	OrgBlock
	Lang     string
	Switches string
	// Headers  string
	// TODO: Change into
	Headers      map[string]string
	ResultsBlock *OrgResultsBlock
}

type OrgResultsBlock struct {
	Hash            string
	ParentCodeBlock *OrgSrcBlock
}

func processHeadline(line string, root OrgRoot) *OrgHeadline {
	matches := root.RegexpMap["headline"].FindStringSubmatch(line)
	headline := new(OrgHeadline)
	headline.RawContent = matches[0]
	headline.Level = len(matches[1])
	headline.State = matches[2]
	headline.Priority = matches[3]
	headline.Title = matches[4]
	headline.Tags = matches[5]
	return headline
}

func processSection(line string, currentOrgElement interface{}, root OrgRoot) interface{} {

	var nextOrgElement interface{}
	if root.RegexpMap["headline"].MatchString(line) {
		headline := processHeadline(line, root)
		nextOrgElement = headline
		return nextOrgElement
	}

	// if already set, then just increment the number of lines
	switch ce := currentOrgElement.(type) {
	case *OrgParagraph:
		ce.RawContent += line + "\n"
		nextOrgElement = ce
	default:
		nextOrgElement = &OrgParagraph{OrgElement: OrgElement{RawContent: line + "\n"}}
	}

	return nextOrgElement
}

func processSrcBlock(content string, lpos int, org *OrgRoot) (interface{}, int) {

	src := &OrgSrcBlock{}

	// Start to scan once again within this code block.
	var srcContentStart int
	var srcContentEnd int
	var srcEnd int

	// position of the next line
	npos := lpos
	for lookahead := lpos; lookahead < len(content); lookahead += 1 {

		if content[lookahead] == LF {
			line := content[npos:lookahead]

			if org.RegexpMap["beginSrc"].MatchString(line) {
				m := org.RegexpMap["beginSrc"].FindStringSubmatchIndex(line)

				inlineBlockInfo := line[m[1]:len(line)]

				// match blockInfo
				infom := org.RegexpMap["blockInfo"].FindStringSubmatch(inlineBlockInfo)
				src.Lang = infom[1]
				src.Switches = infom[2]
				// src.Headers   = infom[3]
				srcContentStart = lpos + len(line) + 1
			}

			if org.RegexpMap["endSrc"].MatchString(line) {
				org.RegexpMap["endSrc"].FindStringSubmatchIndex(line)
				srcContentEnd = npos
				srcEnd = lookahead
				break // done
			}

			// move the position to the end of the line
			npos = lookahead + 1
		}
	}

	// TODO: Actually just `content` not raw content as in order elements
	src.RawContent = content[srcContentStart:srcContentEnd]

	// lookbehind for #+name and #+headers
	// these become concatenated and interpreted later
	src.Headers = make(map[string]string, 0)

	revpos := lpos - 1
	for lookbehind := revpos - 1; lookbehind > 0; lookbehind -= 1 {
		if string(content[lookbehind]) == "\n" {
			line := content[(lookbehind + 1):revpos]

			if org.RegexpMap["blockHeaders"].MatchString(line) {
				headerMatches := org.RegexpMap["blockHeaders"].FindStringSubmatch(line)
				ha := parseHeaderArguments(headerMatches[2] + string(NUL))

				// append to the header args unless already exists
				for k, v := range ha {

					// Do not override even if lookbehind matches
					if _, ok := src.Headers[k]; !ok {
						src.Headers[k] = v
					}
				}

				revpos = lookbehind

			} else if org.RegexpMap["blockName"].MatchString(line) {

				blockNameMatches := org.RegexpMap["blockName"].FindStringSubmatch(line)
				src.Name = blockNameMatches[2]
				break

			} else {

				break
			}
		}
	}

	return src, srcEnd
}

func determineHeadlineHierarchy(cep *OrgHeadline, ne *OrgHeadline) *OrgHeadline {

	var upperParent *OrgHeadline

	// check if parent is root
	if cep.Parent == nil {
		// just use cep then
		ne.Parent = cep
		cep.Elements = append(cep.Elements, ne)
		return cep
	}

	if cep.Level < ne.Level { // subheadline
		ne.Parent = cep
		cep.Elements = append(cep.Elements, ne)
	} else if cep.Level == ne.Level {
		// the grandparent of the siblings
		switch cegp := cep.Parent.(type) {
		case *OrgHeadline:
			ne.Parent = cegp
			cegp.Elements = append(cegp.Elements, ne)

		case *OrgRoot:
			ne.Parent = cegp
			cegp.Elements = append(cegp.Elements, ne)
		default:
			panic("Error during parsing: Could not find correct parent node!")
		}
	} else {
		upperParent = determineHeadlineHierarchy(cep.Parent.(*OrgHeadline), ne)
	}

	return upperParent
}

func parseHeaderArguments(headerArgs string) map[string]string {
	hargs := make(map[string]string, 0)
	var state int
	var keyword string
	var value string
	var pos int

	for lookahead, c := range headerArgs {
		switch c {
		case SPACE, TAB:
			switch state {
			case ORG_HEADER_ARGUMENT_BEGIN:
				// skip
			case ORG_HEADER_ARGUMENT_KEYWORD_BEGIN:
				state = ORG_HEADER_ARGUMENT_KEYWORD_END
				keyword = headerArgs[pos:lookahead]
			case ORG_HEADER_ARGUMENT_KEYWORD_END:
				// skip
			}

		case COLON:
			switch state {
			case ORG_HEADER_ARGUMENT_BEGIN:
				state = ORG_HEADER_ARGUMENT_KEYWORD_BEGIN
				pos = lookahead
			case ORG_HEADER_ARGUMENT_VALUE_BEGIN:
				state = ORG_HEADER_ARGUMENT_KEYWORD_BEGIN
				value = headerArgs[pos : lookahead-1]
				pos = lookahead
				hargs[keyword] = value
			}

			// TODO: match when we are done
		case LF, NUL:
			switch state {
			case ORG_HEADER_ARGUMENT_VALUE_BEGIN:
				value = headerArgs[pos:lookahead]
				hargs[keyword] = value
			}

		default:
			switch state {
			case ORG_HEADER_ARGUMENT_KEYWORD_END:
				state = ORG_HEADER_ARGUMENT_VALUE_BEGIN
				pos = lookahead
			}
		}
	}

	return hargs
}

func Tokenize(content string, options ...*OrgRoot) []interface{} {
	tokens := make([]interface{}, 0)

	// TODO: In case it is nil because we didn't preprocess
	//       we just initialize an OrgRoot.
	var org *OrgRoot
	if len(options) > 0 {
		org = options[0]
	} else {
		org = &OrgRoot{}
		org.RegexpMap = DefaultOrgRegexpMap
		org.Level = 0
	}

	// Initialize the slice of elements
	org.Elements = make([]interface{}, 0)

	// Start with the root as current element
	var currentOrgElement interface{}
	currentOrgElement = org

	// scan using position
	lpos := 0
	line := ""

	for lookahead := 0; lookahead < len(content); lookahead++ {
		// <root>       ::= <section>
		// <section>    ::= <headline> | <fixedwidthblock> | <commentblock> | <affiliatedkeyword> | <block> | <paragraph>
		// <headline>   ::= <headlinestruct> "\n",
		//                  <property-drawer>? <section( without headline of same or higher level)>?,
		//                  <section(without headline of lower level)>
		// <paragraph>  ::= "\n", <section>
		// <block>      ::= <codeblock> | <exampleblock> | <quoteblock> | <verseblock>
		// <codeblock>  ::= <beginsrcblock> | <blockswitches> | <blockheaderargs> | <blockcontent> | <endsrcblock>
		//
		var nextOrgElement interface{}

		// if lpos > lookahead {
		// catchup
		// }

		if content[lookahead] == LF {
			line = content[lpos:lookahead]

			// we can continue having the same element, it is ok

			// Handle empty lines
			// TODO: check for empty lines could be done better
			if len(string(line)) == 0 {
				switch c := currentOrgElement.(type) {
				case *OrgRoot, *OrgHeadline, *OrgParagraph:

					// Just accumulate in case of empty if
					nextOrgElement = &OrgEmptyLine{OrgElement: OrgElement{RawContent: "\n"}}
					goto PUSH_TOKEN

				case *OrgEmptyLine:

					// Continue accumulating the lines and don't change anything
					c.RawContent += "\n"
					nextOrgElement = currentOrgElement
					goto PUSH_TOKEN
				default:
					nextOrgElement = currentOrgElement
				}
			}

			// headline
			if org.RegexpMap["headline"].MatchString(line) {
				headline := processHeadline(line, *org)
				nextOrgElement = headline
				goto PUSH_TOKEN
			}

			// Sections: Must be checked after the headline

			// Match for elements that start with ':'
			// : example

			// Match for elements that start with '#'
			if org.RegexpMap["beginKeywordOrComment"].MatchString(line) {
				// First character in match after '#' is very important!
				m := org.RegexpMap["beginKeywordOrComment"].FindStringSubmatchIndex(line)

				// if no more characters after '#', handle as paragraph
				if string(content[lookahead-1]) == "#" {
					nextOrgElement = &OrgParagraph{OrgElement: OrgElement{RawContent: line + "\n"}}
					goto PUSH_TOKEN
				}

				if string(line[m[1]]) == " " {
					// TODO: This is a Comment token, in case the current element is a token
					// then we have to accumulate it in its RawContent
				}

				rest := line[m[1]:len(line)]
				if org.RegexpMap["beginBlock"].MatchString(rest) {
					// we have a block! Now check its type

					if org.RegexpMap["beginBlock"].MatchString(rest) {
						blockType := org.RegexpMap["beginBlock"].FindStringSubmatch(rest)[1]
						switch blockType {
						case "SRC":
							var endSrcPos int
							nextOrgElement, endSrcPos = processSrcBlock(content, lpos, org)
							lookahead = endSrcPos
							goto PUSH_TOKEN

							// TODO: Implement
							// case "EXAMPLE":
							// case "QUOTE":
							// case "CENTER":
							// case "COMMENT":
							// case "VERSE":
						}
					}
				}
			}

			// If already within a paragraph section, we accumulate the output
			// otherwise we create a new paragraph section
			switch ce := currentOrgElement.(type) {
			case *OrgParagraph:
				ce.RawContent += line + "\n"
				nextOrgElement = ce
				goto PUSH_TOKEN
			}

			// Default is to treat as paragraph
			nextOrgElement = &OrgParagraph{OrgElement: OrgElement{RawContent: line + "\n"}}
			goto PUSH_TOKEN

			//
			// ---------------------------------
			// switch ne := nextOrgElement.(type) {
			// case *OrgSrcBlock:
			// default:
			// }
			// ---------------------------------

		PUSH_TOKEN:
			// store the token, and move position to the beginning of next line
			tokens = append(tokens, nextOrgElement)
			currentOrgElement = nextOrgElement
			lpos = lookahead + 1
		} // content
	} // for lookahead

	return tokens
}

func Parse(tokens []interface{}, currentOrgElement interface{}) interface{} {

	var initialElement interface{}
	initialElement = currentOrgElement

	for _, nextToken := range tokens {

		switch ce := currentOrgElement.(type) {
		case *OrgRoot:

			// (root, paragraph, srcblock) -> headline
			switch ne := nextToken.(type) {
			case *OrgHeadline:
				ce.Elements = append(ce.Elements, ne)
				ne.Parent = ce
				currentOrgElement = ne

			case *OrgParagraph:
				ce.Elements = append(ce.Elements, ne)
				ne.Parent = ce
				currentOrgElement = ne

			default:
				panic("Invalid transition from Root")
			}

		case *OrgHeadline:

			// (headline) -> <child | parent | sibling> headline
			//             | paragraph
			//             | <block>
			switch ne := nextToken.(type) {
			case *OrgHeadline:

				if ce.Level < ne.Level {
					// sub headline
					ce.Elements = append(ce.Elements, ne)
					ne.Parent = ce
					currentOrgElement = ne
				} else {
					// TODO: Non-subheadline cases, look for the proper parent
					switch cep := ce.Parent.(type) {
					case *OrgHeadline: // (headline) -> child,parent, sibling headline
						determineHeadlineHierarchy(cep, ne)
					case *OrgRoot:
						// determineHeadlineHierarchy(cep, ne)
					}
					currentOrgElement = ne
				}
			case *OrgSrcBlock:
				ce.Elements = append(ce.Elements, ne)
				ne.Parent = ce
				currentOrgElement = ne

			case *OrgParagraph:
				ce.Elements = append(ce.Elements, ne)
				ne.Parent = ce
				currentOrgElement = ne
			}

			// section elements
		case *OrgParagraph:
			// (section) -> paragraph | block | headline
			switch ne := nextToken.(type) {
			case *OrgParagraph:
				// It has been accumulated in the scanner already
				// so we just merge them into the same paragraph
				ne.Parent = ce.Parent
				currentOrgElement = ne

			case *OrgSrcBlock:
				switch cep := ce.Parent.(type) {
				case *OrgHeadline: // headline -> (para) -> headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep

				case *OrgRoot: // root ::= para ; para ::= headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep
				}

				currentOrgElement = ne

			case *OrgHeadline: // (para) -> headline
				// <paragraph> ::= <section>
				// <section>   ::= <headline> | <etc...>
				if ce.Parent == nil {
					panic("Error while creating Org mode syntax tree: Headline without parent!")
				}

				switch cep := ce.Parent.(type) {
				case *OrgHeadline: // headline -> (para) -> headline
					determineHeadlineHierarchy(cep, ne)

				case *OrgRoot: // root ::= para ; para ::= headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep
				}
				currentOrgElement = ne

			default:
				// pass
			}

		case *OrgSrcBlock:

			// (section) -> paragraph | block | headline
			switch ne := nextToken.(type) {
			case *OrgParagraph:
				// It has been accumulated in the scanner already
				// so we just merge them into the same paragraph

				switch cep := ce.Parent.(type) {
				case *OrgHeadline: // headline -> (para) -> headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep

				case *OrgRoot: // root ::= para ; para ::= headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep
				}

				currentOrgElement = ne

			case *OrgSrcBlock:
				switch cep := ce.Parent.(type) {
				case *OrgHeadline: // headline -> (para) -> headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep

				case *OrgRoot: // root ::= para ; para ::= headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep
				}

				currentOrgElement = ne

			case *OrgHeadline: // (para) -> headline
				// <paragraph> ::= <section>
				// <section>   ::= <headline> | <etc...>
				if ce.Parent == nil {
					panic("Error while creating Org mode syntax tree: Headline without parent!")
				}

				switch cep := ce.Parent.(type) {
				case *OrgHeadline: // headline -> (para) -> headline
					determineHeadlineHierarchy(cep, ne)

				case *OrgRoot: // root ::= para ; para ::= headline
					cep.Elements = append(cep.Elements, ne)
					ne.Parent = cep
				}
				currentOrgElement = ne

			default:
				// pass
			}
		}
	}

	return initialElement
	// return root // root.Elements is the first section of the tree
}
