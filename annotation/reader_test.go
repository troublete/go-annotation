package annotation

import (
	"testing"

	"github.com/troublete/go-annotation/inspect"
)

func Test_Read(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result, err := Read(inspect.TypeList{
			inspect.Type{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
				Fields: []inspect.Field{
					{
						Comments: []string{
							`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
						},
					},
				},
			},
		}, inspect.FunctionList{
			{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
			},
		})
		if err != nil {
			t.Error(err)
		}
		if result == nil {
			t.Error("expected results")
		}
	})

	t.Run("error funcs", func(t *testing.T) {
		result, err := Read(inspect.TypeList{
			inspect.Type{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
			},
		}, inspect.FunctionList{
			{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"`,
				},
			},
		})
		if err == nil {
			t.Error("expected error")
		}
		if result != nil {
			t.Error("didn't expect results")
		}
	})

	t.Run("error types", func(t *testing.T) {
		result, err := Read(inspect.TypeList{
			inspect.Type{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"`,
				},
			},
		}, inspect.FunctionList{
			{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
			},
		})
		if err == nil {
			t.Error("expected error")
		}
		if result != nil {
			t.Error("didn't expect results")
		}
	})

	t.Run("error fields", func(t *testing.T) {
		result, err := Read(inspect.TypeList{
			inspect.Type{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
				Fields: []inspect.Field{
					{
						Comments: []string{
							`some_annotation{valid="something,",another_valid=sdfsdfwer2342"`,
						},
					},
				},
			},
		}, inspect.FunctionList{
			{
				Comments: []string{
					`some_annotation{valid="something,",another_valid=sdfsdfwer2342"}`,
				},
			},
		})
		if err == nil {
			t.Error("expected error")
		}
		if result != nil {
			t.Error("didn't expect results")
		}
	})
}
