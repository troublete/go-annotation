package analyze

import (
	"fmt"
	"github.com/troublete/go-chariot/inspect"
	"regexp"
	"strings"
	"sync"
)

var (
	DefinitionRe = regexp.MustCompile(`^([a-zA-Z\.\_]+)\s{0,1}\{(.*)\}$`)
	ArgumentRe   = regexp.MustCompile(`^([a-zA-Z\_]+)\=(.*)$`)
	BracketRe    = regexp.MustCompile(`({|})`)
	Separator    = ";"
)

const (
	WarnFormatBrackets       = "bracket structure doesn't match"
	WarnFormatWrongFormat    = "comment '%v' doesn't match the required definition format"
	WarnAttributeWrongFormat = "attribute '%v' doesn't match the required attribute format"
)

type Definition struct {
	Identifier string
	Arguments  map[string]string
}

type DefinitionList []Definition

type Warning string

type Warnings []string

func ExtractDefinitionsForFunction(f inspect.Function, filter *func(bool) bool) (DefinitionList, Warnings) {
	var dlMu sync.Mutex
	var dl DefinitionList
	if f.Comments == nil {
		return dl, nil
	}

	var warnMu sync.Mutex
	var warnings []string

	var wg sync.WaitGroup
	wg.Add(len(f.Comments))
	for _, c := range f.Comments {
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

			pairs := strings.Split(m[2], Separator)
			for _, p := range pairs {
				if p == "" {
					continue
				}

				am := ArgumentRe.FindStringSubmatch(p)

				if len(am) < 3 {
					warnMu.Lock()
					warnings = append(warnings, fmt.Sprintf(WarnAttributeWrongFormat, p))
					warnMu.Unlock()
					return
				}

				def.Arguments[am[1]] = am[2]
			}

			dlMu.Lock()
			defer dlMu.Unlock()
			dl = append(dl, def)
		}(c)
	}
	wg.Wait()

	return dl, warnings
}
