package inspect

// annotation_test{attribute=name,second_attribute=somethingelse}
type TestA struct {
	// annotation.test{attribute=name,another_one=anotheronewith!"§234}
	// annotation.test{csv="alpha,beta,gamma",another_one="test123!"}
	FieldA string
}

// annotation.test{sub=}
func (ta TestA) A() {}

// test.test{something=else}
func (ta *TestA) B() {}

// another.annotation{test=abc123}
func C() {}
