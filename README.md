# gqlgen-apollo-error

ApolloError compliant error function for gqlgen

## Installation

```sh
$ go get -u github.com/s-ichikawa/gqlgen-apollo-error
```

## Quickstart
Return error as:
```go
import (
	...
    "github.com/s-ichikawa/gqlgen-apollo-error"
)

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	...
	if err != nil {
		return nil, gqlgen_apollo_error.UserInputError("Invalid argument value")
	}
	...
}
```
Response:
```json
{
  "errors": [
    {
      "message": "Invalid argument value",
      "path": [
        "todos"
      ],
      "extensions": {
        "code": "BAD_USER_INPUT"
      }
    }
  ],
  "data": null
}
```

## Error Codes
| Error Code | function |
| ------------- | ------------- |
| GRAPHQL_PARSE_FAILED  | SyntaxError  |
| GRAPHQL_VALIDATION_FAILED  | ValidationError  |
| BAD_USER_INPUT  | UserInputError  |
| UNAUTHENTICATED  | AuthenticationError  |
| FORBIDDEN  | ForbiddenError  |
| PERSISTED_QUERY_NOT_FOUND  | PersistedQueryNotFoundError  |
| PERSISTED_QUERY_NOT_SUPPORTED  | PersistedQueryNotSupportedError  |
| INTERNAL_SERVER_ERROR  | None  |

## Extensions
### WithError
Add an error message separately from the `errors.message`:
```go

```
Response
```json
{
  "errors": [
    {
      "message": "Invalid argument value",
      "path": [
        "todos"
      ],
      "extensions": {
        "code": "BAD_USER_INPUT",
        "exception": {
          "error": "something wrong"
        }
      }
    }
  ],
  "data": null
}
```

To set the StackTrace to be displayed, set as follows before returning an error:
```go
gqlgen_apollo_error.SetError(gqlgen_apollo_error.Config{StackTrace: true})
```

```json
{
  "errors": [
    {
      "message": "Invalid argument value",
      "path": [
        "todos"
      ],
      "extensions": {
        "code": "BAD_USER_INPUT",
        "exception": {
          "error": "something wrong",
          "stacktrace": [
            "gqlgen-todo/graph.(*queryResolver).Todos",
            "    /path/to/gqlgen-todo/graph/schema.resolvers.go:28",
            "...more lines...",
          ]
        }
      }
    }
  ],
  "data": null
}
```

### WithValue
Add any set of Key and Value
```go
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	...
	if err != nil {
		return nil, gqlgen_apollo_error.UserInputError(
			"Invalid argument value", 
			gqlgen_apollo_error.WithValue("argumentName", "id"),
		)
	}
```
Response
```json
{
  "errors": [
    {
      "message": "Unknown Error",
      "path": [
        "todos"
      ],
      "extensions": {
        "argumentName": "id",
        "code": "BAD_USER_INPUT"
      }
    }
  ],
  "data": null
}
```