package djiwifi

func must[T any](in T, err error) T {
	if err != nil {
		panic(err)
	}
	return in
}

func cannotFail(err error) {
	if err != nil {
		panic(err)
	}
}
