package errors

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var (
	ErrGenericBadReq = &ValidationError{Message: ""}
	ErrIDRequired    = &ValidationError{Message: ""} // Attribute ID is required
	ErrInvalidID     = &ValidationError{Message: ""} // ID must be an integer
	ErrIDPositive    = &ValidationError{Message: ""} // ID must be a positive integer
)
