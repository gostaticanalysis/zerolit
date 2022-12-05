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
