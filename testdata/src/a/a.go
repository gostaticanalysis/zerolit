package a

func f1() int {
	return 0 // OK
}

func f2() struct{} {
	v := struct{}{}
	return v // want "zero value should return as a literal"
}

func f3() (v struct{}) {
	return v // want "zero value should return as a literal"
}

func f4() struct{ int } {
	v := struct{ int }{100}
	return v // OK
}

func f5() struct{ int } {
	var v struct{ int }
	v = struct{ int }{100}
	return v // OK
}

func f6() int {
	var v int
	v = 0
	return v // want "zero value should return as a literal"
}
