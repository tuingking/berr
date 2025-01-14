# Berr

Package `berr` is an error handling library

`go get github.com/tuingking/berr`

<!-- toc -->

- [Berr](#berr)
  - [Using berr](#using-berr)
    - [Wrapping errors](#wrapping-errors)
    - [Creating errors](#creating-errors)

<!-- tocstop -->

## Using berr

### Wrapping errors

`berr.Wrap` adds context to an error while preserving the original error.

```golang
err := Cook(usr, "burger")
if err != nil {
  // wrap the error if you want to add more context
  return nil, berr.Wrap(err, "error cooking burger"))
}
```

`berr.Wrapx` like `berr.Wrap` but without message.

```golang
err := Cook(usr, "burger")
if err != nil {
  // wrap the error if you want to add more context
  return nil, berr.Wrapx(err))
}
```

### Creating errors

Creating errors is simple via `berr.New`.

```golang
func (req *Request) Validate() error {
  if req.ID == "" {
    // or return a new error at the source if you prefer
    return berr.New("id required")
  }
  return nil
}
```

`berr.Newx` like `berr.New` but without stacktrace.

```golang
var ErrInternalServer = berr.Newx("error internal server")
```