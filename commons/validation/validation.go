package validation

import (
	"strings"
)

type format struct {
	Value    string
	Msg      string
	Param    string
	Location string
}

//Validate type for running validations on data in incoming requests.
type Validate struct {
	ValidationResult []format
}

//ValidateEmail validates email
func (e *Validate) ValidateEmail(email string, message string) {
	f := format{Value: email, Msg: message, Param: "email", Location: "body"}
	trimmed := strings.Trim(email, " ")
	index := strings.Index(trimmed, ".")
	if len(trimmed)-4 != index {
		e.ValidationResult = append(e.ValidationResult, f)
		return
	}
	chars := trimmed[index+1:]
	shouldContainAt := trimmed[:index]
	if (strings.Compare(chars, "com") != 0) ||
		(strings.Index(shouldContainAt, "@") <= 0) ||
		(strings.Count(shouldContainAt, "@") != 1) {
		e.ValidationResult = append(e.ValidationResult, f)
		return
	}

}

//ValidatePasswordLength validates password
func (e *Validate) ValidatePasswordLength(password string, min int, max int, message string) {
	f := format{Value: password, Msg: message, Param: "password", Location: "body"}
	if !(min < len(password) && max > len(password)) {
		e.ValidationResult = append(e.ValidationResult, f)

	}

}

//IsPassword does ..
func (e *Validate) IsPassword(password string, message string) {
	f := format{Value: password, Msg: message, Param: "password", Location: "body"}
	if password == "" {
		e.ValidationResult = append(e.ValidationResult, f)
	}
}

//SerializeErrors does...
func (e *Validate) SerializeErrors() struct {
	Errors []struct {
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
} {
	serialized := make([]struct {
		Message string `json:"message"`
		Field   string `json:"Field,omitempty" bson:"Field,omitempty"`
	}, len(e.ValidationResult))
	for i, v := range e.ValidationResult {
		serialized[i] = struct {
			Message string `json:"message"`
			Field   string `json:"Field,omitempty" bson:"Field,omitempty"`
		}{v.Msg, v.Param}
	}
	d := struct {
		Errors []struct {
			Message string `json:"message"`
			Field   string `json:"field,omitempty"`
		} `json:"errors"`
	}{}

	d.Errors = []struct {
		Message string "json:\"message\""
		Field   string "json:\"field,omitempty\""
	}(serialized)

	return d
}
