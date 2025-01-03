# forms

[![Go Reference](https://pkg.go.dev/badge/cattlecloud.net/go/forms.svg)](https://pkg.go.dev/cattlecloud.net/go/forms)
[![License](https://img.shields.io/github/license/cattlecloud/forms?color=7C00D8&style=flat-square&label=License)](https://github.com/cattlecloud/forms/blob/main/LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/cattlecloud/forms/ci.yaml?style=flat-square&color=0FAA07&label=Tests)](https://github.com/cattlecloud/forms/actions/workflows/ci.yaml)

`forms` provides a way to parse http Request forms using a schema

### Getting Started

The `forms` package can be added to a project by running:

```shell
go get cattlecloud.net/go/forms@latest
```

```go
import "cattlecloud.net/go/forms"
```

### Examples

##### parsing http request

```go
var (
  name string
  age  int
)

err := forms.Parse(request, forms.Schema{
  "NAME": forms.String(&name),
  "AGE":  forms.Int(&age),
})
```

##### about requests

Typically the HTTP request will be given to you in the form of an http handler,
e.g.

```go
func(w http.ResponseWriter, r *http.Request) {
  _ = r.ParseForm()
  // now r form data is available to parse
}
```

### License

The `cattlecloud.net/go/forms` module is open source under the [BSD](LICENSE) license.
