package analyze

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	Separator    = ","
	DefinitionRe = regexp.MustCompile(`^([a-zA-Z\.\_]+)\s{0,1}\{(.*)\}$`)
	ArgumentRe   = regexp.MustCompile(
		fmt.Sprintf(
			`(?P<key>[a-zA-Z_]+)(=("(?P<value>[^"]*)"|(?P<value>[^%s]*)))?%s?`,
			Separator,
			Separator,
		),
	)
	BracketRe    = regexp.MustCompile(`({|})`)
	AssignmentRe = regexp.MustCompile(`=`)
)

const (
	WarnFormatBrackets       = "bracket structure doesn't match"
	WarnFormatWrongFormat    = "comment '%v' doesn't match the required definition format"
	WarnAttributeWrongFormat = "attributes '%v' don't match fully the required attribute format"

	TrueString = "TRUE"
)

type Definition struct {
	Identifier string            `json:"identifier"`
	Arguments  map[string]string `json:"arguments"`
}

type DefinitionList []Definition

type Warning string

type Warnings []string

type Spec interface {
	Doc() []string
}

// ExtractDefinitionsOnSpec extracts protocol matching annotations from a functions documentation
// this is used to build a structure of annotations assigned to functions
func ExtractDefinitionsOnSpec(s Spec, filter *func(bool) bool) (DefinitionList, Warnings) {
	var dlMu sync.Mutex
	var dl DefinitionList
	if s.Doc() == nil {
		return dl, nil
	}

	var warnMu sync.Mutex
	var warnings []string

	var wg sync.WaitGroup
	wg.Add(len(s.Doc()))
	for _, c := range s.Doc() {
		go func(c string) {
			defer wg.Done()

			nbrackets := len(BracketRe.FindAllStringSubmatch(c, -1))
			if (nbrackets % 2) != 0 {
				warnMu.Lock()
				warnings = append(warnings, WarnFormatBrackets)
				warnMu.Unlock()
				return
			}

			if filter != nil {
				if (*filter)(DefinitionRe.MatchString(c)) == false {
					return // skip filtered comments
				}
			}

			def := Definition{}
			m := DefinitionRe.FindStringSubmatch(c)

			if len(m) < 3 {
				warnMu.Lock()
				warnings = append(warnings, fmt.Sprintf(WarnFormatWrongFormat, c))
				warnMu.Unlock()
				return
			}

			def.Identifier = m[1]
			def.Arguments = map[string]string{}

			attributeString := m[2]
			matches := ArgumentRe.FindAllStringSubmatch(m[2], -1)
			for _, m := range matches {
				attributeString = strings.Replace(attributeString, m[0], "", 1) // remove matched part

				key := m[ArgumentRe.SubexpIndex("key")]
				value := m[4]
				if value == "" {
					value = m[5]
				}
				if value == "" && !AssignmentRe.MatchString(m[0]) {
					value = TrueString
				}

				def.Arguments[key] = value
			}

			if len(attributeString) > 0 {
				warnMu.Lock()
				warnings = append(warnings, fmt.Sprintf(WarnAttributeWrongFormat, m[2]))
				warnMu.Unlock()
				return
			}

			dlMu.Lock()
			defer dlMu.Unlock()
			dl = append(dl, def)
		}(c)
	}
	wg.Wait()

	return dl, warnings
}
