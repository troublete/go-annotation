package inspect

import (
	"strings"
	"testing"
)

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
