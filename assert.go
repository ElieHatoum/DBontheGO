package dbonthego

func check(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}
