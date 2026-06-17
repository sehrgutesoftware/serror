package serror

import "fmt"

// CodeResolver is a function that resolves an error code to its error value.
type CodeResolver func(code string) (error, bool)

// Hydrate restores an error chain from its serializable format.
//
// It recursively restores the error and its cause chain from a tree structure
// of errors. The provided [CodeResolver] function is used to restore the
// original error value from the error code.
func Hydrate(node *Tree, resolve CodeResolver) {
	if node == nil {
		return
	}

	if node.Code != "" {
		if resolved, ok := resolve(node.Code); ok {
			node.err = resolved
		} else {
			node.err = fmt.Errorf("unresolved error (code %s): %s", node.Code, node.Message)
		}
	}

	for _, child := range node.Children {
		Hydrate(child, resolve)
	}
}
