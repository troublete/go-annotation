package inspect

import (
	"fmt"
	"strings"
	"testing"
)

func Test_FunctionDoc(t *testing.T) {
	comments := []string{
		"Test",
	}

	f := Function{
		Comments: comments,
	}

	if strings.Join(comments, "") != strings.Join(f.Doc(), "") {
		t.Error("expected document returned be same as comments")
	}
}

func Test_FunctionListFind(t *testing.T) {
	t.Run("find successful", func(t *testing.T) {
		var fl FunctionList
		fl = append(fl, Function{
			Name: "A",
		})

		f := fl.Find("A")
		if f == nil {
			t.Error("expected 'A' to be found")
		}
	})

	t.Run("find unsuccessful", func(t *testing.T) {
		var fl FunctionList
		fl = append(fl, Function{
			Name: "B",
		})

		f := fl.Find("A")
		if f != nil {
			t.Error("expected 'A' not to be found")
		}
	})
}

func Test_TypeListFind(t *testing.T) {
	t.Run("find successful", func(t *testing.T) {
		var tl TypeList
		tl = append(tl, Type{
			Name: "A",
		})

		tt := tl.Find("A")
		if tt == nil {
			t.Error("expected 'A' to be found")
		}
	})

	t.Run("find unsuccessful", func(t *testing.T) {
		var tl TypeList
		tl = append(tl, Type{
			Name: "B",
		})

		tt := tl.Find("A")
		if tt != nil {
			t.Error("expected 'A' not to be found")
		}
	})
}

func Test_FindAllFunctions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		funcs, err := FindAllFunctions("./internal/success")
		if err != nil {
			t.Error(err)
		}

		if len(funcs) != 7 {
			t.Errorf("wanted 7, got %v", len(funcs))
		}

		ta := funcs.Find("TestA")
		if ta == nil {
			t.Error("failed to find 'TestA'")
		}

		if ta.Receiver != nil ||
			ta.Package != "success" ||
			len(ta.Comments) != 3 ||
			strings.Join(ta.Comments, "") != "Commentsanother lineand another" ||
			ta.Name != "TestA" ||
			ta.FilePath != "internal/success/a.go" {
			t.Error("TestA failed expectation")
		}

		aa := funcs.Find("AA")
		if aa == nil {
			t.Error("failed to find 'AA'")
		}

		if aa.Receiver == nil ||
			aa.Receiver.Pointer == true ||
			aa.Receiver.ReceiverType != "TestTypeA" ||
			aa.Package != "success" ||
			aa.Comments != nil ||
			aa.Name != "AA" ||
			aa.FilePath != "internal/success/a.go" {
			t.Error("AA failed expectation")
		}

		ab := funcs.Find("AB")
		if ab == nil {
			t.Error("failed to find 'AB'")
		}

		if ab.Receiver == nil ||
			ab.Receiver.Pointer == false ||
			ab.Receiver.ReceiverType != "TestTypeA" ||
			ab.Package != "success" ||
			ab.Comments != nil ||
			ab.Name != "AB" ||
			ab.FilePath != "internal/success/a.go" {
			t.Error("AB failed expectation")
		}

		cf := funcs.Find("ComplicatedFunction")
		if cf == nil {
			t.Error("failed to find 'ComplicatedFunction'")
		}

		if cf.Receiver != nil ||
			cf.Package != "success" ||
			cf.Comments != nil ||
			cf.Name != "ComplicatedFunction" ||
			cf.FilePath != "internal/success/a.go" {
			t.Error("ComplicatedFunction failed expectation")
		}

		tb := funcs.Find("TestB")
		if tb == nil {
			t.Error("failed to find 'TestB'")
		}

		if tb.Receiver != nil ||
			tb.Package != "sub" ||
			tb.Comments != nil ||
			tb.Name != "TestB" ||
			tb.FilePath != "internal/success/sub/b.go" {
			t.Error("TestB failed expectation")
		}

		ba := funcs.Find("BA")
		if ba == nil {
			t.Error("failed to find 'BA'")
		}

		if ba.Receiver == nil ||
			ba.Receiver.Pointer != false ||
			ba.Receiver.ReceiverType != "TestTypeB" ||
			ba.Package != "sub" ||
			ba.Comments != nil ||
			ba.Name != "BA" ||
			ba.FilePath != "internal/success/sub/b.go" {
			t.Error("BA failed expectation")
		}

		bb := funcs.Find("BB")
		if bb == nil {
			t.Error("failed to find 'BB'")
		}

		if bb.Receiver == nil ||
			bb.Receiver.Pointer == false ||
			bb.Receiver.ReceiverType != "TestTypeB" ||
			bb.Package != "sub" ||
			len(bb.Comments) != 2 ||
			bb.Name != "BB" ||
			bb.FilePath != "internal/success/sub/b.go" ||
			strings.Join(bb.Comments, "") != "test commentanother line" {
			t.Error("BB failed expectation")
		}
	})
	t.Run("error", func(t *testing.T) {
		funcs, err := FindAllFunctions("./internal/error")
		if funcs != nil {
			t.Error("expected failed parsing")
		}

		if err == nil {
			t.Error("expected error")
		}
	})
}

func Test_FindAllTypes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		types, err := FindAllTypes("./internal/success")
		if types == nil {
			t.Error("expected failed parsing")
		}

		if err != nil {
			t.Error("expected error")
		}

		if len(types[0].Comments) > 0 ||
			types[0].FilePath != "internal/success/sub/b.go" ||
			types[0].Name != "TestTypeB" ||
			types[0].Package != "sub" ||
			len(types[0].Fields) != 0 {
			t.Error("TestTypeB failed expectation")
		}

		if strings.Join(types[1].Comments, "") != "Test comment on TestTypeA" ||
			types[1].FilePath != "internal/success/a.go" ||
			types[1].Name != "TestTypeA" ||
			types[1].Package != "success" ||
			len(types[1].Fields) != 5 ||
			strings.Join(types[1].Fields[0].Comments, "") != "Test comment on ValueA" ||
			types[1].Fields[0].Name != "ValueA" ||
			types[1].Fields[0].Type.String() != "string" ||
			types[1].Fields[0].Tags["literal"] != "tag" ||
			types[1].Fields[0].Tags["json"] != "something" ||
			strings.Join(types[1].Fields[1].Comments, "") != "" ||
			types[1].Fields[1].Name != "ComplexType" ||
			types[1].Fields[1].Type.String() != "bytes.Buffer" ||
			types[1].Fields[1].Tags != nil ||
			strings.Join(types[1].Fields[2].Comments, "") != "" ||
			types[1].Fields[2].Name != "ComplexTypePointer" ||
			types[1].Fields[2].Type.String() != "*bytes.Buffer" ||
			types[1].Fields[2].Tags != nil ||
			strings.Join(types[1].Fields[3].Comments, "") != "" ||
			types[1].Fields[3].Name != "StringPointer" ||
			types[1].Fields[3].Type.String() != "*LocalType" ||
			types[1].Fields[3].Tags != nil ||
			types[1].Fields[3].Type.PackageNameImplied != true ||
			types[1].Fields[3].Type.Package != "success" ||
			strings.Join(types[1].Fields[4].Comments, "") != "" ||
			types[1].Fields[4].Name != "String" ||
			types[1].Fields[4].Type.String() != "LocalType" ||
			types[1].Fields[4].Tags != nil ||
			types[1].Fields[4].Type.PackageNameImplied != true ||
			types[1].Fields[4].Type.Package != "success" {
			fmt.Println(
				strings.Join(types[1].Comments, ""), strings.Join(types[1].Comments, "") != "Test comment on TestTypeA", "\n",
				types[1].FilePath, types[1].FilePath != "internal/success/a.go", "\n",
				types[1].Name, types[1].Name != "TestTypeA", "\n",
				types[1].Package, types[1].Package != "success", "\n",
				len(types[1].Fields), len(types[1].Fields) != 4, "\n",
				strings.Join(types[1].Fields[0].Comments, ""), strings.Join(types[1].Fields[0].Comments, "") != "Test comment on ValueA", "\n",
				types[1].Fields[0].Name, types[1].Fields[0].Name != "ValueA", "\n",
				types[1].Fields[0].Type.String(), types[1].Fields[0].Type.String() != "string", "\n",
				types[1].Fields[0].Tags["literal"], types[1].Fields[0].Tags["literal"] != "tag", "\n",
				types[1].Fields[0].Tags["json"], types[1].Fields[0].Tags["json"] != "something", "\n",
				strings.Join(types[1].Fields[1].Comments, ""), strings.Join(types[1].Fields[1].Comments, "") != "", "\n",
				types[1].Fields[1].Name, types[1].Fields[1].Name != "ComplexType", "\n",
				types[1].Fields[1].Type.String(), types[1].Fields[1].Type.String() != "bytes.Buffer", "\n",
				types[1].Fields[1].Tags, types[1].Fields[1].Tags != nil, "\n",
				strings.Join(types[1].Fields[2].Comments, ""), strings.Join(types[1].Fields[2].Comments, "") != "", "\n",
				types[1].Fields[2].Name, types[1].Fields[2].Name != "ComplexTypePointer", "\n",
				types[1].Fields[2].Type.String(), types[1].Fields[2].Type.String() != "*bytes.Buffer", "\n",
				types[1].Fields[2].Tags, types[1].Fields[2].Tags != nil, "\n",
				strings.Join(types[1].Fields[3].Comments, ""), strings.Join(types[1].Fields[3].Comments, "") != "", "\n",
				types[1].Fields[3].Name, types[1].Fields[3].Name != "StringPointer", "\n",
				types[1].Fields[3].Type.String(), types[1].Fields[3].Type.String() != "*LocalType", "\n",
				types[1].Fields[3].Tags, types[1].Fields[3].Tags != nil, "\n",
				types[1].Fields[3].Type.PackageNameImplied, types[1].Fields[3].Type.PackageNameImplied != true, "\n",
				types[1].Fields[3].Type.Package, types[1].Fields[3].Type.Package != "success", "\n",
				strings.Join(types[1].Fields[4].Comments, ""), strings.Join(types[1].Fields[4].Comments, "") != "",
				types[1].Fields[4].Name, types[1].Fields[4].Name != "String",
				types[1].Fields[4].Type.String(), types[1].Fields[4].Type.String() != "LocalType",
				types[1].Fields[4].Tags, types[1].Fields[4].Tags != nil,
				types[1].Fields[4].Type.PackageNameImplied, types[1].Fields[4].Type.PackageNameImplied != true,
				types[1].Fields[4].Type.Package, types[1].Fields[4].Type.Package != "success",
			)
			t.Error("TestTypeA failed expectation")
		}
	})

	t.Run("error", func(t *testing.T) {
		types, err := FindAllTypes("./internal/error")
		if types != nil {
			t.Error("expected failed parsing")
		}

		if err == nil {
			t.Error("expected error")
		}
	})
}

func Test_FieldDoc(t *testing.T) {
	comments := []string{
		"Test",
	}

	f := Field{
		Comments: comments,
	}

	if strings.Join(comments, "") != strings.Join(f.Doc(), "") {
		t.Error("expected document returned be same as comments")
	}
}

func Test_TypeDoc(t *testing.T) {
	comments := []string{
		"Test",
	}

	tt := Type{
		Comments: comments,
	}

	if strings.Join(comments, "") != strings.Join(tt.Doc(), "") {
		t.Error("expected document returned be same as comments")
	}
}
