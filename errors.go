package errors

import (
	"fmt"
)

type CtxKey string

const (
	CtxKeyParent = "parent_err"
)

// ContextualError extends basic error interface with two value-setter and value-getter.
type ContextualError interface {
	error

	// Populates error with a new value for its context.
	//
	// The provided key must be comparable and should not be of type
	// string or any other built-in type to avoid collisions between
	// packages using ctx-errors. Users of WithValue should define their own
	// types for keys. To avoid allocating when assigning to an
	// interface{}, context keys often have concrete type
	// struct{}. Alternatively, exported context key variables' static
	// type should be a pointer or interface.
	WithValue(key, value interface{}) ContextualError

	// Value returns the value associated with this context for key, or nil
	// if no value is associated with key. Successive calls to Value with
	// the same key returns the same result.
	//
	// Use context values only for request-scoped data that transits
	// processes and API boundaries, not for passing optional parameters to
	// functions.
	//
	// A key identifies a specific value in a Context. Functions that wish
	// to store values in Context typically allocate a key in a global
	// variable then use that key as the argument to context.WithValue and
	// Context.Value. A key can be any type that supports equality;
	// packages should define keys as an unexported type to avoid
	// collisions.
	//
	// Packages that define a Context key should provide type-safe accessors
	// for the values stored using that key:
	Value(key interface{}) interface{}
}

type ctxErr struct {
	msg     string
	context map[interface{}]interface{}
}

// See "error" interface for details.
func (err *ctxErr) Error() string {
	return fmt.Sprintf("%s: with context: %+v", err.msg, err.context)
}

// See ContextualError.WithValue() for details.
func (err *ctxErr) WithValue(key, value interface{}) ContextualError {
	if key == nil {
		panic("nil key")
	}
	if err.context == nil {
		err.context = map[interface{}]interface{}{}
	}
	err.context[key] = value
	return err
}

// See ContextualError.Value() for details.
func (err *ctxErr) Value(key interface{}) interface{} {
	if value, ok := err.context[key]; ok {
		return value
	}
	if parent, ok := err.context[CtxKeyParent]; ok {
		if ctxParent, ok := parent.(ContextualError); ok {
			return ctxParent.Value(key)
		}
	}
	return nil
}

// Produces new contextual error with empty context.
func New(msg string) ContextualError {
	return &ctxErr{
		msg: msg,
	}
}

// Produces new formatted contextual error with empty context.
func Errorf(format string, args ...interface{}) ContextualError {
	return New(fmt.Sprintf(format, args...))
}

// Wraps existing error with a contextualized one (with empty context).
func Wrap(err error, msg string) ContextualError {
	newErr := &ctxErr{msg: msg}
	return newErr.WithValue(CtxKeyParent, err)
}

// Wraps existing error with a formatted contextualized one (with empty context).
func Wrapf(err error, format string, args ...interface{}) ContextualError {
	return Wrap(err, fmt.Sprintf(format, args...))
}

// Populates error with a new value for its context.
// If error is not contextual - it will be wrapped.
//
// See ContextualError.WithValue() for details.
func WithValue(err error, key, value interface{}) ContextualError {
	if ctxErr, ok := err.(ContextualError); ok {
		return ctxErr.WithValue(key, value)
	}
	return Wrap(err, "error").WithValue(key, value)
}

// See ContextualError.Value() for details.
func Value(err error, key interface{}) interface{} {
	if ctxErr, ok := err.(ContextualError); ok {
		return ctxErr.Value(key)
	}
	return nil
}
