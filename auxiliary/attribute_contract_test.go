package auxiliary

import "testing"

func Test_AttributeContract(t *testing.T) {
	for _, tc := range []struct {
		a      Attributes
		p      Preconditions
		result []string
		pass   bool
	}{
		{
			a: Attributes{
				"test":         "a",
				"passthrough":  "whatever",
				"filtered_out": "something",
			},
			p: Preconditions{
				"test": func(v string) bool {
					return v == "a"
				},
				"passthrough": PreconditionAlwaysAllow,
			},
			result: []string{"test", "passthrough"},
			pass:   true,
		},
		{
			a: Attributes{
				"test":  "a",
				"fails": "critical_value",
			},
			p: Preconditions{
				"test": PreconditionAlwaysAllow,
				"fails": func(_ string) bool {
					return false
				},
			},
			result: []string{"test"},
			pass:   false,
		},
		{
			a: Attributes{
				"test": "a",
			},
			p: Preconditions{
				"test": PreconditionAlwaysAllow,
				"not_set": func(_ string) bool {
					return true
				},
			},
			result: []string{"test"},
			pass:   false,
		},
	} {
		r := AttributeContract(tc.a, tc.p)
		for k, _ := range r.Attributes {
			unexpectedResult := true
			for _, v := range tc.result {
				if v == k {
					unexpectedResult = false
				}
			}
			if unexpectedResult {
				t.Error("unexpected attribute")
			}
		}
		if r.Valid != tc.pass {
			t.Error("contract failed")
		}

	}
}
