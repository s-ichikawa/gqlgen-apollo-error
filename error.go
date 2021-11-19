package gqlgen_apollo_error

import (
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"runtime"
)

type ErrorCode string

const (
	SyntaxErrorCode                     ErrorCode = "GRAPHQL_PARSE_FAILED"
	ValidationErrorCode                 ErrorCode = "GRAPHQL_VALIDATION_FAILED"
	UserInputErrorCode                  ErrorCode = "BAD_USER_INPUT"
	AuthenticationErrorCode             ErrorCode = "UNAUTHENTICATED"
	ForbiddenErrorCode                  ErrorCode = "FORBIDDEN"
	PersistedQueryNotFoundErrorCode     ErrorCode = "PERSISTED_QUERY_NOT_FOUND"
	PersistedQueryNotSupportedErrorCode ErrorCode = "PERSISTED_QUERY_NOT_SUPPORTED"
	NoneCode                            ErrorCode = "INTERNAL_SERVER_ERROR"
)

var aerr = &apolloError{}

type apolloError struct {
	stackTrace bool
}

type Config struct {
	StackTrace bool
}

func SetError(c Config) {
	aerr = &apolloError{
		stackTrace: c.StackTrace,
	}
}

type Extensions map[string]interface{}
type Extension interface {
	apply(extensions Extensions)
}

type extensionFunc func(Extensions)

func (f extensionFunc) apply(extensions Extensions) {
	f(extensions)
}

func WithError(err error) Extension {
	return extensionFunc(func(extensions Extensions) {
		exception := map[string]interface{}{}
		exception["error"] = err.Error()
		if aerr.stackTrace {
			var st []string
			i := 4
			for {
				pt, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				f := runtime.FuncForPC(pt).Name()
				st = append(st, f)
				st = append(st, fmt.Sprintf("    %s:%d", file, line))
				i++
			}
			exception["stacktrace"] = st
		}
		extensions["exception"] = exception
	})
}

func WithValue(key string, value interface{}) Extension {
	return extensionFunc(func(extensions Extensions) {
		extensions[key] = value
	})
}

func (e *apolloError) NewErrorf(code ErrorCode, msg string, options ...Extension) *gqlerror.Error {
	extentions := map[string]interface{}{
		"code": code,
	}

	for _, o := range options {
		o.apply(extentions)
	}

	return &gqlerror.Error{
		Message:    msg,
		Extensions: extentions,
	}
}

// SyntaxError is the GraphQL operation string contains a syntax error.
func SyntaxError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(SyntaxErrorCode, msg, options...)
}

// ValidationError is the GraphQL operation is not valid against the server's schema.
func ValidationError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(ValidationErrorCode, msg, options...)
}

// UserInputError is the GraphQL operation includes an invalid value for a field argument.
func UserInputError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(UserInputErrorCode, msg, options...)
}

// AuthenticationError is the server failed to authenticate with a required data source, such as a REST API.
func AuthenticationError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(AuthenticationErrorCode, msg, options...)
}

// ForbiddenError is the server was unauthorized to access a required data source, such as a REST API.
func ForbiddenError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(ForbiddenErrorCode, msg, options...)
}

// PersistedQueryNotFoundError is a client sent the hash of a query string to execute
// via automatic persisted queries, but the query was not in the APQ cache.
func PersistedQueryNotFoundError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(PersistedQueryNotFoundErrorCode, msg, options...)
}

// PersistedQueryNotSupportedError is A client sent the hash of a query string
// to execute via automatic persisted queries, but the server has disabled APQ.
func PersistedQueryNotSupportedError(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(PersistedQueryNotSupportedErrorCode, msg, options...)
}

// None is An unspecified error occurred.
func None(msg string, options ...Extension) *gqlerror.Error {
	return aerr.NewErrorf(NoneCode, msg, options...)
}
