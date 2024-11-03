package annotation

import (
	"fmt"
	"strings"

	"github.com/troublete/go-annotation/analyze"
	"github.com/troublete/go-annotation/inspect"
)

type AnnotatedFunction struct {
	Function    inspect.Function       `json:"function"`
	Annotations analyze.DefinitionList `json:"annotations,omitempty"`
}

type AnnotatedType struct {
	Type        inspect.Type           `json:"type"`
	Annotations analyze.DefinitionList `json:"annotations,omitempty"`
	Fields      []AnnotatedField       `json:"fields"`
}

type AnnotatedField struct {
	Field       inspect.Field          `json:"field"`
	Annotations analyze.DefinitionList `json:"annotations,omitempty"`
}

type Result struct {
	Functions []AnnotatedFunction `json:"functions,omitempty"`
	Types     []AnnotatedType     `json:"types,omitempty"`
}

func Read(types inspect.TypeList, funcs inspect.FunctionList) (*Result, error) {
	result := Result{}
	for _, t := range types {
		at := AnnotatedType{
			Type: t,
		}

		defs, warnings := analyze.ExtractDefinitionsOnSpec(t, analyze.FilterCommentNoAnnotation())
		if len(warnings) > 0 {
			return nil, fmt.Errorf("warnings occured: %v", strings.Join(warnings, ", "))
		}
		at.Annotations = defs

		for _, f := range t.Fields {
			af := AnnotatedField{
				Field: f,
			}

			defs, warnings := analyze.ExtractDefinitionsOnSpec(f, analyze.FilterCommentNoAnnotation())
			if len(warnings) > 0 {
				return nil, fmt.Errorf("warnings occured: %v", strings.Join(warnings, ", "))
			}
			af.Annotations = defs
			at.Fields = append(at.Fields, af)
		}

		result.Types = append(result.Types, at)
	}

	for _, f := range funcs {
		af := AnnotatedFunction{
			Function: f,
		}
		defs, warnings := analyze.ExtractDefinitionsOnSpec(f, analyze.FilterCommentNoAnnotation())
		if len(warnings) > 0 {
			return nil, fmt.Errorf("warnings occured: %v", strings.Join(warnings, ", "))
		}
		af.Annotations = defs

		result.Functions = append(result.Functions, af)
	}

	return &result, nil
}
