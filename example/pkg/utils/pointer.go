package utils

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
