//go:build exclude

package sub

func TestB() {}

type TestTypeB struct{}

func (tta TestTypeB) BA() {}

// test comment
// another line
func (tta *TestTypeB) BB() {}
