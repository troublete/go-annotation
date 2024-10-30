package analyze

import (
	"fmt"
	"github.com/troublete/go-chariot/inspect"
	"testing"
)

func Test_ExtractDefinitionsForFunction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		for _, tc := range []struct {
			c string
			d Definition
		}{
			{
				`chariot.route{some_key_without_value=;something_else={"json":"example"}}`,
				Definition{
					Identifier: "chariot.route",
					Arguments: map[string]string{
						"some_key_without_value": "",
						"something_else":         `{"json":"example"}`,
					},
				},
			},
			{
				`user.custom{test=attribute;test_number=123;}`,
				Definition{
					Identifier: "user.custom",
					Arguments: map[string]string{
						"test":        "attribute",
						"test_number": "123",
					},
				},
			},
		} {
			t.Run(tc.c, func(t *testing.T) {
				defs, warnings := ExtractDefinitionsForFunction(inspect.Function{
					Comments: []string{
						tc.c,
					},
				}, nil)

				if len(warnings) > 0 {
					t.Error("expected no warnings")
				}

				if defs[0].Identifier != tc.d.Identifier {
					t.Errorf("expected identifier to match, but didn't (has=%v, want=%v)", defs[0].Identifier, tc.d.Identifier)
				}

				for k, v := range defs[0].Arguments {
					if av, ok := tc.d.Arguments[k]; !ok || v != av {
						t.Errorf("argument didn't match (has=%v, want=%v)", k+":"+v, k+":"+av)
					}
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		for _, tc := range []struct {
			c string
			w Warnings
		}{
			{
				`chariot.route{some_key_without_value=;something_else={"json":"example"}`, // missing end
				[]string{
					WarnFormatBrackets,
				},
			},
			{
				`chariot route{some_key_without_value=;something_else={"json":"example"}}`, // wrong ident
				[]string{
					fmt.Sprintf(WarnFormatWrongFormat, `chariot route{some_key_without_value=;something_else={"json":"example"}}`),
				},
			},
			{
				`chariot_route{some_key_without_value;something_else={"json":"example"}}`, // wrong attribute
				[]string{
					fmt.Sprintf(WarnAttributeWrongFormat, "some_key_without_value"),
				},
			},
		} {
			t.Run(tc.c, func(t *testing.T) {
				defs, warnings := ExtractDefinitionsForFunction(inspect.Function{
					Comments: []string{
						tc.c,
					},
				}, nil)

				if len(warnings) < 1 {
					t.Error("expected warnings")
				}

				if len(defs) > 0 {
					t.Error("didn't expect definitions")
				}

				for idx, w := range warnings {
					if w != tc.w[idx] {
						t.Errorf("warning didn't match, got\n'%v'\n, want \n'%v'", w, tc.w[idx])
					}
				}
			})
		}
	})

	t.Run("empty comments", func(t *testing.T) {
		defs, err := ExtractDefinitionsForFunction(inspect.Function{}, nil)
		if err != nil {
			t.Error("didn't expect error")
		}

		if len(defs) > 0 {
			t.Error("expected no definitions")
		}
	})

	t.Run("test filter", func(t *testing.T) {
		defs, err := ExtractDefinitionsForFunction(inspect.Function{
			Comments: []string{
				"comment not matching format",
			},
		}, FilterCommentNotInForm())
		if err != nil {
			t.Error("didn't expect error")
		}

		if len(defs) > 0 {
			t.Error("expected no definitions")
		}
	})
}
