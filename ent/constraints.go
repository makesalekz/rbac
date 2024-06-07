package ent

func NewConstraintError(msg string, wrap error) *ConstraintError {
	return &ConstraintError{msg, wrap}
}

func NewNotFoundError(label string) *NotFoundError {
	return &NotFoundError{label}
}
