package djictl

func cannotFail(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](in T, err error) T {
	if err != nil {
		panic(err)
	}
	return in
}
