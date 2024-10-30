//go:build exclude

package success

// Comments
// another line
// and another
func TestA() {}

type TestTypeA struct{}

func (ttb TestTypeA) AA()  {}
func (ttb *TestTypeA) AB() {}

func ComplicatedFunction() {
	a := func() {}
	a()
}
