package _error

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
)

// Op represents the operation or function where an error occurred.
type Op string

// Detail provides additional information about the error.
// used in 'error.detail' json response.
type Detail string

// Title represents a short summary of the error.
// used in 'message' json response.
type Title string

// Kind represents the category of the error.
// used in 'error.kind' json response.
type Kind uint8

// Error kinds representing various types of errors.
const (
	Other          Kind = iota // Unexpected error
	InvalidRequest             // Error caused by an invalid request
	Exist                      // Error when a resource already exists
	NotExist                   // Error when a resource does not exist
	Validation                 // Error caused by validation failure
	Forbidden                  // Error due to insufficient permissions
	Database                   // Database-related error
	Internal                   // Internal server error
)

// String returns the string representation of an error Kind.
func (k Kind) String() string {
	switch k {
	case Other:
		return "other_error"
	case InvalidRequest:
		return "invalid_request"
	case Exist:
		return "resource_already_exist"
	case NotExist:
		return "resource_not_found"
	case Validation:
		return "validation_error"
	case Forbidden:
		return "forbidden"
	case Database:
		return "database_error"
	case Internal:
		return "internal_server_error"
	default:
		return "unknown_error_kind"
	}
}

// ProblemParameter represents a parameter associated with a specific problem.
type ProblemParameter struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type Params []ProblemParameter

// TODO: refactor all code using this type of error 'Problem'.

// Problem represents detailed information about an error.
type Problem struct {
	Op     Op     // Operation where the error occurred
	Kind   Kind   // Category of the error
	Title  Title  // Short summary of the error
	Detail Detail // Additional details about the error
	Params Params // Parameters associated with the error
	Err    error  // Nested or wrapped error
}

// E constructs a new Problem instance with the provided arguments.
//
// The function can accept various types of arguments, including Op, Kind, Title,
// Detail, and error. It combines them to create a comprehensive error instance.
//
// Parameters:
//   - args: Variadic arguments to define the problem.
//
// Returns:
//   - error: A new Problem instance.
//
// Panics:
//   - If called with no arguments.
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to _error.E with no arguments")
	}

	e := &Problem{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case Title:
			e.Title = arg
		case Detail:
			e.Detail = arg
		case string:
			e.Err = errors.New(arg)
		case *Problem:
			errCopy := *arg
			e.Err = &errCopy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("_error.E: bad call from %s:%d %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	prev, ok := e.Err.(*Problem)
	if !ok {
		return e
	}

	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	if prev.Title == e.Title {
		prev.Title = ""
	}

	if e.Title == "" && prev.Title != "" {
		e.Title = prev.Title
	}

	if prev.Detail == e.Detail {
		prev.Detail = ""
	}

	if e.Detail == "" {
		e.Detail = prev.Detail
		prev.Detail = ""
	}

	return e
}

func (p *Problem) Unwrap() error {
	return p.Err
}

func (p *Problem) Error() string {
	return p.Err.Error()
}

// OpStack generates a stack of operations from a chain of Problem instances.
//
// Parameters:
//   - err: The error instance to analyze.
//
// Returns:
//   - []string: A slice of operation names in reverse order of occurrence.
func OpStack(err error) []string {
	type o struct {
		Op    string
		Order int
	}

	e := err
	i := 0
	var os []o

	for errors.Unwrap(e) != nil {
		var errProblem *Problem
		if errors.As(e, &errProblem) {
			if errProblem.Op != "" {
				op := o{Op: string(errProblem.Op), Order: i}
				os = append(os, op)
			}
		}
		e = errors.Unwrap(e)
		i++
	}

	sort.Slice(os, func(i, j int) bool { return os[i].Order > os[j].Order })

	var ops []string
	for _, op := range os {
		ops = append(ops, op.Op)
	}

	return ops
}
