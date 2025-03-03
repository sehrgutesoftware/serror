package serror

import "errors"

// Tree structure representing an error and its cause chain.
type Tree struct {
	Code     string
	err      error
	Children []*Tree
}

// Error returns the error message.
func (e *Tree) Error() string {
	return e.err.Error()
}

// Is checks if the error is the same as the target.
func (e *Tree) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	if e.err == target {
		return true
	}

	if target, ok := target.(*Tree); ok {
		if e.Code == target.Code {
			return true
		}
		if e.err == target.err {
			return true
		}
	}

	if coder, ok := target.(errorCoder); ok && e.Code == coder.ErrorCode() {
		return true
	}

	return false
}

// As checks if the error can be converted to the target.
func (e *Tree) As(target any) bool {
	return errors.As(e.err, target)
}

// Unwrap returns the underlying errors.
func (e *Tree) Unwrap() []error {
	errs := make([]error, len(e.Children))
	for i, child := range e.Children {
		errs[i] = child
	}
	return errs
}
