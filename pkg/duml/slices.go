package duml

func array1ToSlice[E any, T [1]E](in T) []E {
	return []E{in[0]}
}
