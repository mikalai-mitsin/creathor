package errs

func NewEquipmentNotFound() *Error {
	return &Error{
		Code:    ErrorCodeNotFound,
		Message: "Equipment not found.",
		Params:  map[string]string{},
	}
}
