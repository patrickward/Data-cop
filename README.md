# Datacop

[![Go Reference](https://pkg.go.dev/badge/github.com/patrickward/datacop.svg)](https://pkg.go.dev/github.com/patrickward/datacop)
[![Go Report Card](https://goreportcard.com/badge/github.com/patrickward/datacop)](https://goreportcard.com/report/github.com/patrickward/datacop)

Datacop is a validation library for Go applications that prioritizes type safety and readability over magic.

## Features

- 🔒 Type-safe validations with generics support
- ⚡ No reflection-based struct tags
- 🎯 Explicit, readable validation logic
- 🔗 Chainable validations
- 🌳 Support for nested validation groups
- 🎭 Conditional validations
- ✨ Common set of built-in validation functions

## Installation

```bash
go get github.com/patrickward/datacop
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/patrickward/datacop"
)

func main() {
    v := datacop.New()
    
    // Simple field validation
    username := "john"
    v.Check(is.Required(username), "username", "username is required")
    
    // Chainable validation
    email := "invalid-email"
    v.Field("email", email).
        Check(is.Required(email), "email is required").
        Check(is.Email(email), "invalid email format")
    
    if v.HasErrors() {
        fmt.Println(v.Errors())
        return
    }
}
```

## Design Philosophy

Datacop intentionally favors explicit validation over struct tag-based validation for:
- Better compile-time type safety
- Clearer validation logic
- More flexible conditional validations
- Easier testing
- No reflection overhead

## Documentation

See the [package documentation](https://pkg.go.dev/github.com/patrickward/datacop) for detailed usage examples and API reference.

### Common Validation Patterns

```go
v := datacop.New()

// Explicit checks 
v.Check(is.Required(username), "username is required")
v.Check(is.Email(email), "invalid email format")
v.CheckStandalone(is.Min(18)(age), "must be 18 or older")

// Password validation
v.Field("password", password).
    Check(is.Required(password), "password is required").
    Check(is.MinLength(8)(password), "password too short").
    Check(is.Match(`[A-Z]`)(password), "must contain uppercase").
    Check(is.Match(`[a-z]`)(password), "must contain lowercase").
    Check(is.Match(`[0-9]`)(password), "must contain number")

// Grouped validation
userGroup := v.Group("user")
userGroup.Field("name", name).
    Check(is.Required(name), "name is required")
userGroup.Field("age", age).
    Check(is.Min(18)(age), "must be 18 or older")

// Conditional validation
v.Field("company", company).
    When(isEmployed).
    Check(is.Required(company), "company required")
```
