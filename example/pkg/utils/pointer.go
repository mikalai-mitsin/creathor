package utils

import con "golang.org/x/exp/constraints"

func Pointer[T any](v T) *T {
	return &v
}

func Value[T any](v *T) T {
	var d T
	if v != nil {
		return *v
	}
	return d
}

func ToAnySlice[T any](input []T) []any {
	b := make([]any, len(input))
	for i := range input {
		b[i] = input[i]
	}
	return b
}

func ChangeType[T con.Integer, F con.Integer](from []F) []T {
	res := make([]T, len(from))
	for i, e := range from {
		res[i] = T(e)
	}
	return res
}
