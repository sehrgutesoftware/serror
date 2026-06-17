package serror

// errorCoder is an interface that allows an error to uniquely identify itself
type errorCoder interface {
	// ErrorCode returns a unique error code for the error. This code is used
	// to serialize the error in a way that allows it to be restored to its
	// original type. Returning an empty string means that the error does
	// not have a unique code.
	ErrorCode() string
}

// Encode transforms an error chain into a serializable format.
//
// It recursively unwraps the error and its cause chain to create a tree
// structure of errors. If an error implements the errorCoder interface,
// the error code is included in the output structure.
func Encode(e error) *Tree {
	if e == nil {
		return nil
	}

	node := &Tree{
		err:      e,
		Children: make([]*Tree, 0),
		Message:  e.Error(),
	}

	if ec, ok := e.(errorCoder); ok {
		node.Code = ec.ErrorCode()
	}

	switch e := e.(type) {
	case interface{ Unwrap() error }:
		if child := Encode(e.Unwrap()); child != nil {
			node.Children = append(node.Children, child)
		}
	case interface{ Unwrap() []error }:
		for _, err := range e.Unwrap() {
			if child := Encode(err); child != nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node
}
