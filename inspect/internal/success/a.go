//go:build exclude

package success

import "bytes"

// Comments
// another line
// and another
func TestA() {}

// Test comment on TestTypeA
type TestTypeA struct {
	// Test comment on ValueA
	ValueA             string `literal:tag,json:something`
	ComplexType        bytes.Buffer
	ComplexTypePointer *bytes.Buffer
	StringPointer      *string
}

func (ttb TestTypeA) AA()  {}
func (ttb *TestTypeA) AB() {}

func ComplicatedFunction() {
	a := func() {}
	a()
}

// test{}
type TestTypeNotMatching string

// test{}
type TestTypeNotMatchingPointer *string

// test{}
type TestTypeNotMatchingTypePointer *bytes.Buffer
