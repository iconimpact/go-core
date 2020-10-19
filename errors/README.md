# Errors

Package errors provides simple error handling.

It allows to set the kind of error this is, a HTTP message and record the strack trace of error.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconmobile-dev/go-core/errors
```

## Usage and Examples

`errors.E` builds an error value from its arguments.

    There must be at least one argument or E panics.
    The type of each argument determines its meaning.
    If more than one argument of a given type is presented,
    only the last one is recorded.

    The types are:
        string
            The HTTP message for the API user.
        errors.Kind
            The class of error, such as permission failure.
        error
            The underlying error that triggered this one.

    If the error is printed, only those items that have been
    set to non-zero values will appear in the result.

    If Kind is not specified or Other, we set it to the Kind of
    the underlying error.

```go
errors.E(errors.Unprocessable, "HTTP response message")
```

```go
err := db.Get(&obj, query)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return errors.E(err, errors.NotFound)
    }

    return errors.E(err, errors.Internal)
}
```

`errors.ToHTTPResponse` creates a string to be used for HTTP response by chaining the underlying application errors HTTPMessage.

```go
err := errors.E(errors.Unprocessable, "HTTP response message 1")
err = errors.E(err, "HTTP response message 2")

errors.ToHTTPResponse(err)
response: "HTTP response message 1: HTTP response message 2"
```

<br>

## Kinds of errors

When the error is used in a router the `Kind` is usually mapped to a HTTP response code.

```go
Other                     // Unclassified error
Unauthorized              // Unauthorized (401)
Forbidden                 // Forbidden (403)
NotFound                  // Not found (404)
Conflict                  // Conflict (409)
Unprocessable             // Unprocessable, invalid request data (422)
Internal                  // Internal server error (500)
```