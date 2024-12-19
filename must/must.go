package must

func NotFail(err error) {
	if err != nil {
		panic(err)
	}
}

func BeTrue(v bool) {
	if !v {
		panic("expect a boolean true value")
	}
}
