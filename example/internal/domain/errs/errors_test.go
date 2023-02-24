package errs

import (
	"reflect"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

func TestError_AddParam(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Code:    0,
				Message: "",
				Params:  map[string]string{},
			},
			args: args{
				key:   "Betty",
				value: "Piptik",
			},
			want: &Error{
				Code:    0,
				Message: "",
				Params: map[string]string{
					"Betty": "Piptik",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			e.AddParam(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("Is() = %v, want %v", e, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Code:    16,
				Message: "This is error.",
				Params: map[string]string{
					"first":  "foo",
					"second": "bar",
				},
			},
			want: "{\"code\":16,\"message\":\"This is error.\",\"params\":{\"first\":\"foo\",\"second\":\"bar\"}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	type args struct {
		tgt error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok true",
			fields: fields{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "bar",
				},
			},
			args: args{
				tgt: NewUnexpectedBehaviorError("bar"),
			},
			want: true,
		},
		{
			name: "ok false",
			fields: fields{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "no bar",
				},
			},
			args: args{
				tgt: NewUnexpectedBehaviorError("bar"),
			},
			want: false,
		},
		{
			name: "ok false",
			fields: fields{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "no bar",
				},
			},
			args: args{
				tgt: errors.New("bar"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			if got := e.Is(tt.args.tgt); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_SetCode(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	type args struct {
		code ErrorCode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Code:    0,
				Message: "",
				Params:  nil,
			},
			args: args{
				code: 16,
			},
			want: &Error{
				Code:    16,
				Message: "",
				Params:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			e.SetCode(tt.args.code)
			if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("Is() = %v, want %v", e, tt.want)
			}
		})
	}
}

func TestError_SetMessage(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Code:    0,
				Message: "",
				Params:  nil,
			},
			args: args{
				message: "this is message!",
			},
			want: &Error{
				Code:    0,
				Message: "this is message!",
				Params:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			e.SetMessage(tt.args.message)
			if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("Is() = %v, want %v", e, tt.want)
			}
		})
	}
}

func TestError_SetParams(t *testing.T) {
	type fields struct {
		Code    ErrorCode
		Message string
		Params  map[string]string
	}
	type args struct {
		params map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name: "ok",
			fields: fields{
				Code:    0,
				Message: "",
				Params:  nil,
			},
			args: args{
				params: map[string]string{
					"new": "key",
					"ad":  "off",
				},
			},
			want: &Error{
				Code:    0,
				Message: "",
				Params: map[string]string{
					"new": "key",
					"ad":  "off",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Params:  tt.fields.Params,
			}
			e.SetParams(tt.args.params)
			if !reflect.DeepEqual(e, tt.want) {
				t.Errorf("Is() = %v, want %v", e, tt.want)
			}
		})
	}
}

func TestFromValidationError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "ok validation error object",
			args: args{
				err: validation.NewError("dsa", "asd"),
			},
			want: &Error{
				Code:    3,
				Message: "asd",
				Params:  map[string]string{},
			},
		},
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "no validation",
			args: args{
				err: errors.New("no validation"),
			},
			want: nil,
		},
		{
			name: "ok validation errors domain error",
			args: args{
				err: validation.Errors{
					"first": &Error{
						Code:    16,
						Message: "Text of error",
						Params:  nil,
					},
				},
			},
			want: &Error{
				Code:    3,
				Message: "The form sent is not valid, please correct the errors below.",
				Params: map[string]string{
					"first": "Text of error",
				},
			},
		},
		{
			name: "ok validation errors error object",
			args: args{
				err: validation.Errors{
					"first": &validation.ErrorObject{},
				},
			},
			want: &Error{
				Code:    3,
				Message: "The form sent is not valid, please correct the errors below.",
				Params: map[string]string{
					"first": "",
				},
			},
		},
		{
			name: "ok validation errors",
			args: args{
				err: validation.Errors{
					"first":  validation.NewError("key 1", "value 1"),
					"second": validation.NewError("key 2", "value 2"),
				},
			},
			want: &Error{
				Code:    3,
				Message: "The form sent is not valid, please correct the errors below.",
				Params:  map[string]string{"first": "value 1", "second": "value 2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromValidationError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		code    ErrorCode
		message string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "ok",
			args: args{
				code:    16,
				message: "New error",
			},
			want: &Error{
				Code:    16,
				Message: "New error",
				Params:  map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInvalidFormError(t *testing.T) {
	tests := []struct {
		name string
		want *Error
	}{
		{
			name: "ok",
			want: NewError(
				ErrorCodeInvalidArgument,
				"The form sent is not valid, please correct the errors below.",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInvalidFormError(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInvalidFormError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUnexpectedBehaviorError(t *testing.T) {
	type args struct {
		details string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "ok",
			args: args{
				details: "test details",
			},
			want: &Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "test details",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnexpectedBehaviorError(tt.args.details); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUnexpectedBehaviorError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_renderErrorMessage(t *testing.T) {
	type args struct {
		object validation.ErrorObject
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: func() args {
				obj := validation.ErrorObject{}.
					SetCode("12").
					SetMessage("simple message {{.first}} {{.second}}").
					SetParams(map[string]interface{}{
						"first":  "foo",
						"second": "bar",
					})
				return args{
					object: obj.(validation.ErrorObject),
				}
			}(),
			want: "simple message foo bar",
		},
		{
			name: "bad message",
			args: func() args {
				obj := validation.ErrorObject{}.SetCode("12").
					SetMessage("{{ .text | asd }}").
					SetParams(map[string]interface{}{
						"first":  "foo",
						"second": "bar",
					})
				return args{
					object: obj.(validation.ErrorObject),
				}
			}(),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderErrorMessage(tt.args.object); got != tt.want {
				t.Errorf("renderErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
