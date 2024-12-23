package datacop

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StandaloneErrorKey is the key used for standalone errors, i.e. global errors
const StandaloneErrorKey = "__standalone__"

type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Validator struct {
	errors map[string][]ValidationError
}

// New creates a new validator instance
//
// Example usage:
func New() *Validator {
	return &Validator{
		errors: make(map[string][]ValidationError),
	}
}

// CheckStandalone performs a standalone validation and adds an error if it fails
func (v *Validator) CheckStandalone(valid bool, message string) bool {
	if !valid {
		v.AddStandaloneError(message)
		return false
	}
	return true
}

// Check performs a field validation and adds an error if it fails
func (v *Validator) Check(valid bool, field, message string) bool {
	if !valid {
		v.AddError(field, message)
		return false
	}
	return true
}

// Error implements the error interface
func (v *Validator) Error() string {
	parts := make([]string, 0, len(v.errors))

	if standalone := v.StandaloneErrors(); len(standalone) > 0 {
		parts = append(parts, fmt.Sprintf("global: [%s]", strings.Join(standalone, ", ")))
	}

	for field, errs := range v.errors {
		if field == StandaloneErrorKey {
			continue
		}
		if len(errs) > 0 {
			messages := make([]string, len(errs))
			for i, err := range errs {
				messages[i] = err.Message
			}
			parts = append(parts, fmt.Sprintf("%s: [%s]", field, strings.Join(messages, ", ")))
		}
	}

	return strings.Join(parts, " | ")
}

// AddStandaloneError adds a standalone error
func (v *Validator) AddStandaloneError(message string) {
	v.AddError(StandaloneErrorKey, message)
}

// AddError adds an error for a specific field
func (v *Validator) AddError(field, message string) {
	// Ensure the current validator is initialized
	if v.errors == nil {
		v.errors = make(map[string][]ValidationError)
	}

	v.errors[field] = append(v.errors[field], ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasStandaloneErrors returns true if there are any standalone errors
func (v *Validator) HasStandaloneErrors() bool {
	return v.HasErrorFor(StandaloneErrorKey)
}

// HasErrors returns true if there are any validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// ErrorFor returns the string error message for a field
func (v *Validator) ErrorFor(field string) string {
	if errs, exists := v.errors[field]; exists && len(errs) > 0 {
		messages := make([]string, len(errs))
		for i, err := range errs {
			messages[i] = err.Message
		}
		return strings.Join(messages, ", ")
	}
	return ""
}

// HasErrorFor returns true if the field has any errors for a field
func (v *Validator) HasErrorFor(field string) bool {
	errs, exists := v.errors[field]
	return exists && len(errs) > 0
}

// StandaloneErrors returns all standalone error messages
func (v *Validator) StandaloneErrors() []string {
	if errs, exists := v.errors[StandaloneErrorKey]; exists {
		messages := make([]string, len(errs))
		for i, err := range errs {
			messages[i] = err.Message
		}
		return messages
	}
	return nil
}

// Errors returns a map of field names and their string error messages
func (v *Validator) Errors() map[string]string {
	fields := make(map[string]string)
	for field, errs := range v.errors {
		if len(errs) > 0 {
			fields[field] = v.ErrorFor(field)
		}
	}
	return fields
}

// ValidationErrors returns all validation errors as a map of field names to their errors
func (v *Validator) ValidationErrors() map[string][]ValidationError {
	return v.errors
}

// Merge combines another validator's errors into this one. The other validator is not modified.
func (v *Validator) Merge(other *Validator) {
	// Ensure the current validator is initialized
	if v.errors == nil {
		v.errors = make(map[string][]ValidationError)
	}

	for field, errs := range other.errors {
		v.errors[field] = append(v.errors[field], errs...)
	}
}

// MarshalJSON implements json.Marshaler for the Validator type
func (v *Validator) MarshalJSON() ([]byte, error) {
	fields := make(map[string]string)

	for field, errs := range v.errors {
		if len(errs) > 0 {
			fields[field] = v.ErrorFor(field)
		}
	}

	return json.Marshal(struct {
		Errors map[string]string `json:"fields,omitempty"`
	}{
		Errors: fields,
	})
}

// Clear removes all errors from the validator instance
func (v *Validator) Clear() {
	v.errors = make(map[string][]ValidationError)
}

// FieldValidation enables chain validation for a specific field
type FieldValidation struct {
	field string
	value any
	v     *Validator
}

// Field starts a validation chain for the given field
//
// Example usage:
// v := datacop.New()
// v.Field("username", username).
//
//	Check(Required(username), "username is required").
//	Check(MinLength(3)(username), "username must be at least 3 characters")
//	Check(MaxLength(255)(username), "username must be at most 255 characters")
//
//	if v.HasErrorFor("username") {
//	     fmt.Println(v.ErrorFor("username"))
//	}
func (v *Validator) Field(name string, value any) *FieldValidation {
	return &FieldValidation{
		field: name,
		value: value,
		v:     v,
	}
}

// Check performs a validation in the chain
//
// Example usage:
// v := datacop.New()
// v.Field("username", username).
//
//	Check(Required(username), "username is required").
//	Check(MinLength(3)(username), "username must be at least 3 characters")
func (f *FieldValidation) Check(valid bool, message string) *FieldValidation {
	f.v.Check(valid, f.field, message)
	return f
}

// Group represents a group of related validations
type Group struct {
	name string
	v    *Validator
}

// When represents a conditional validation
type When struct {
	condition bool
	field     string
	value     any
	v         *Validator
}

// Group starts a validation group
func (v *Validator) Group(name string) *Group {
	return &Group{name: name, v: v}
}

// Field starts a validation chain for the given field in the group
func (g *Group) Field(name string, value any) *FieldValidation {
	fullName := g.name + "." + name
	return g.v.Field(fullName, value)
}

// When starts a conditional validation
func (f *FieldValidation) When(condition bool) *When {
	return &When{
		condition: condition,
		field:     f.field,
		value:     f.value,
		v:         f.v,
	}
}

// Check performs a validation in the chain
func (w *When) Check(valid bool, message string) *When {
	if w.condition {
		w.v.Check(valid, w.field, message)
	}
	return w
}

// When performs a validation in the chain
func (w *When) When(condition bool) *When {
	w.condition = w.condition && condition
	return w
}
