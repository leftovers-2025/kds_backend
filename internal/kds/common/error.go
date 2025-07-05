package common

type ValidationError struct {
	message string
}

func NewValidationError(err error) error {
	return &ValidationError{
		message: err.Error(),
	}
}

func (e *ValidationError) Error() string {
	return e.message
}

type InternalError struct {
	message string
}

func NewInternalError(err error) error {
	return &InternalError{
		message: err.Error(),
	}
}

func (e *InternalError) Error() string {
	return e.message
}

type NotFoundError struct {
	message string
}

func NewNotFoundError(err error) error {
	return &NotFoundError{
		message: err.Error(),
	}
}

func (e *NotFoundError) Error() string {
	return e.message
}
