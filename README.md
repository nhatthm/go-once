# Once

A simple library to do something exactly once.

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-once)](https://github.com/nhatthm/go-once/releases/latest)
[![Build Status](https://github.com/nhatthm/go-once/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-once/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-once/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-once)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/once)](https://goreportcard.com/report/go.nhat.io/once)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/once)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

TBD

## Prerequisites

- `Go >= 1.21`

## Install

```bash
go get go.nhat.io/once
```

## Types

| Type           | Description                                                         |
|----------------|---------------------------------------------------------------------|
| `Once`         | An alias of [`sync.Once`](https://pkg.go.dev/sync#Once).            |
| `Func`         | A function that can be executed exactly once.                       |
| `Value`        | A function that returns value can be executed exactly once.         |
| `Values`       | A function that returns values can be executed exactly once.        |
| `FuncMap`      | A map of functions that can be executed exactly once.               |
| `ValueMap`     | A map of functions that return value can be executed exactly once.  |
| `ValuesMap`    | A map of functions that return values can be executed exactly once. |
| `LazyValueMap` | A map contains values that are initialized on the first access.     |

## Examples

```go
package once_test

import (
    "fmt"

    "go.nhat.io/once"
)

func ExampleLazyValueMap() {
    type Person struct {
        ID   string
        Name string
    }

    people := once.LazyValueMap[string, *Person]{
        New: func(key string) *Person {
            return &Person{ID: key}
        },
    }

    instance1 := people.Get("1")
    instance2 := people.Get("1")

    fmt.Println(instance2.Name)

    instance1.Name = "John Doe"

    fmt.Println(instance2.Name)

    // Output:
    //
    // John Doe
}

```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
