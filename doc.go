/*
Package datacop provides a fluent validation library for Go applications.

Note: The validator is not thread-safe. Each validator instance should be used within a single goroutine, or external synchronization should be used when accessing from multiple goroutines.

# Design Philosophy

This package intentionally favors explicit validation over struct tag-based validation for several reasons:
- Better compile-time type safety and IDE support
- Clearer validation logic that's easier to read and maintain
- More flexible support for complex conditional validations
- Easier to unit test individual validation rules
- No reflection overhead

The package offers two main validation styles:
1. Direct validation using Check methods
2. Chainable validation using Field builder

# Basic Usage

The simplest way to use the validator is with direct Check calls:

	v := datacop.New()

	// Validate username
	if !v.Check(Required(username), "username", "username is required") {
		return v.Errors() // returns map[string]string
	}

	// Validate multiple fields
	v.Check(Required(email), "email", "email is required")
	v.Check(Email(email), "email", "invalid email format")
	v.Check(MinLength(8)(password), "password", "password too short")

	if v.HasErrors() {
		return v.Error() // returns formatted error string
	}

# Chainable Validation

For more complex validations, use the fluent Field builder:

	v := datacop.New()

	// Chain validations for username
	v.Field("username", username).
		Check(is.Required(username), "username is required").
		Check(is.MinLength(3)(username), "username too short").
		Check(is.MaxLength(255)(username), "username too long")

	// Chain validations for password
	v.Field("password", password).
		Check(is.Required(password), "password is required").
		Check(is.MinLength(8)(password), "password too short").
		Check(is.Match(`[A-Z]`)(password), "must contain uppercase").
		Check(is.Match(`[0-9]`)(password), "must contain number")

# Grouped Validation

For nested structures, use Group to namespace validations:

	v := datacop.New()

	// Validate user details group
	userGroup := v.Group("user")
	userGroup.Field("name", name).
		Check(is.Required(name), "name is required")
	userGroup.Field("age", age).
		Check(is.Min(18)(age), "must be 18 or older")

	// Validate address group
	addrGroup := v.Group("address")
	addrGroup.Field("street", street).
		Check(is.Required(street), "street is required")
	addrGroup.Field("city", city).
		Check(is.Required(city), "city is required")

# Conditional Validation

Conditional validations using When are evaluated sequentially. When a condition is false, all subsequent checks are skipped until the next When condition:

	v := datacop.New()

	// Check is skipped because When(false)
	v.Field("company", company).
		When(false).
		Check(is.Required(company), "company required")

	// Multiple conditions
	v.Field("role", role).
		Check(is.Required(role), "role is required").     // Always runs
		When(isAdmin).                                 // Only if isAdmin is true
		Check(is.In("admin", "super")(), "invalid role"). //   will this check run
		When(hasPermission).                           // Only if hasPermission is true
		Check(is.NotZero(), "permission required")        //   will this check run

	// When conditions are combined with AND logic
	v.Field("department", dept).
		When(isEmployee).
		When(isFullTime).
		Check(is.Required(dept), "department required")    // Runs only if isEmployee AND isFullTime

Note: Each When condition affects only the Check calls that follow it, until another When is encountered. The validation chain is processed sequentially from left to right.

# Standalone Errors

For validations not tied to specific fields:

	v := datacop.New()

	// Add standalone error
	if !isValid {
		v.AddStandaloneError("general validation failed")
	}

	// Check standalone condition
	v.CheckStandalone(password == confirmPassword, "passwords do not match")

# Custom Validation Functions

Creating custom validation functions is straightforward - any function that returns a bool can be used:

	// Simple custom validation
	func IsEven(value any) bool {
	    v, ok := value.(int)
	    if !ok {
	        return false
	    }
	    return v%2 == 0
	}

	// Usage
	v.Check(IsEven(age), "age", "age must be even")

	// Custom validation with parameters
	func MultipleOf(n int) ValidationFunc {
	    return func(value any) bool {
	        v, ok := value.(int)
	        if !ok {
	            return false
	        }
	        return v%n == 0
	    }
	}

	// Usage
	v.Check(MultipleOf(3)(age), "age", "age must be multiple of 3")

# Common Validation Functions

The package provides many built-in validation functions:

		// Basic validations
		is.Required(value)              // checks if value is non-empty
	 	is.NotZero(value)               // checks if numeric value is not zero
		is.Match(`[a-zA-Z0-9]+`)(value) // regex pattern

		// Comparison validations

		is.Between(1, 100)(value)      // numeric range
		is.Equal(10)(value)            // equal to value
		is.EqualStrings("a", "b")      // equal strings
		is.In("a", "b", "c")(value)   // value in set
		is.AllIn("a", "b", "c")([]string{"a", "b"}) // all values in set
		is.NoDuplicates()([]string{})  // unique values in slice
		is.Min(18)(value)              // minimum value
		is.Max(65)(value)              // maximum value
		is.MinLength(5)(value)         // minimum length
		is.MaxLength(10)(value)        // maximum length
		is.GreaterThan(100)(value)     // greater than value
		is.LessThan(100)(value)        // less than value
		is.GreaterOrEqual(100)(value)  // greater or equal to value
		is.LessOrEqual(100)(value)     // less or equal to value

		// Time-based validations

		is.Before(time.Now())          // time before now
		is.After(time.Now())           // time after now
		is.BetweenTime(start, end)     // time between two values

		// String regex validations
		is.Email(value)                // email format
		is.Phone(value)                // phone number format

# Error Handling

Multiple ways to access validation errors:

	v.HasErrors()               // returns true if any errors exist
	v.HasErrorFor("field")      // checks for field-specific errors
	v.ErrorFor("field")         // gets error message for field
	v.Errors()                  // returns map[string]string of all errors
	v.Error()                   // returns formatted error string
	v.ValidationErrors()        // returns full error structs
	v.StandaloneErrors()        // returns non-field-specific errors

# Common Patterns

Password validation example:

	v := datacop.New()
	v.Field("password", password).
		Check(is.Required(password), "password is required").
		Check(is.MinLength(8)(password), "password too short").
		Check(is.Match(`[A-Z]`)(password), "must contain uppercase").
		Check(is.Match(`[a-z]`)(password), "must contain lowercase").
		Check(is.Match(`[0-9]`)(password), "must contain number")

Form validation example:

	type Form struct {
		Username string
		Email    string
		Age      int
	}

	func ValidateForm(form Form) error {
		v := datacop.New()

		v.Field("username", form.Username).
			Check(is.Required(form.Username), "username is required").
			Check(is.MinLength(3)(form.Username), "username too short")

		v.Field("email", form.Email).
			Check(is.Required(form.Email), "email is required").
			Check(is.Email(form.Email), "invalid email format")

		v.Field("age", form.Age).
			Check(is.Min(18)(form.Age), "must be 18 or older")

		if v.HasErrors() {
			return v
		}
		return nil
	}
*/
package datacop
