package datacop_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patrickward/datacop"
	"github.com/patrickward/datacop/is"
)

func TestValidator_BasicValidation(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         any
		validationFn  func(any) bool
		message       string
		expectError   bool
		expectedError string
	}{
		{
			name:         "valid required string",
			field:        "username",
			value:        "johndoe",
			validationFn: is.Required,
			message:      "username is required",
			expectError:  false,
		},
		{
			name:          "invalid required string",
			field:         "username",
			value:         "",
			validationFn:  is.Required,
			message:       "username is required",
			expectError:   true,
			expectedError: "username is required",
		},
		{
			name:         "valid email",
			field:        "email",
			value:        "test@example.com",
			validationFn: is.Email,
			message:      "invalid email format",
			expectError:  false,
		},
		{
			name:          "invalid email",
			field:         "email",
			value:         "not-an-email",
			validationFn:  is.Email,
			message:       "invalid email format",
			expectError:   true,
			expectedError: "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := datacop.New()
			v.Check(tt.validationFn(tt.value), tt.field, tt.message)

			assert.Equal(t, tt.expectError, v.HasErrors())
			if tt.expectError {
				assert.Equal(t, tt.expectedError, v.ErrorFor(tt.field))
			}
		})
	}
}

func TestValidator_MultipleErrors(t *testing.T) {
	v := datacop.New()

	// Add multiple errors for the same field
	v.Check(false, "password", "too short")
	v.Check(false, "password", "needs uppercase")
	v.Check(false, "password", "needs number")

	assert.True(t, v.HasErrors())
	assert.True(t, v.HasErrorFor("password"))
	assert.Equal(t, "too short, needs uppercase, needs number", v.ErrorFor("password"))

	// Test multiple fields
	v.Check(false, "username", "too short")
	assert.True(t, v.HasErrorFor("username"))

	// The Error() method should contain all errors
	errorStr := v.Error()
	assert.Contains(t, errorStr, "password: [too short, needs uppercase, needs number]")
	assert.Contains(t, errorStr, "username: [too short]")
}

func TestValidator_StandaloneErrors(t *testing.T) {
	v := datacop.New()

	// Add standalone error
	v.CheckStandalone(false, "database connection failed")
	assert.True(t, v.HasStandaloneErrors())

	// Add field error
	v.Check(false, "username", "invalid username")

	// Both types of errors should be present
	assert.True(t, v.HasErrors())
	errorStr := v.Error()
	assert.Contains(t, errorStr, "database connection failed")
	assert.Contains(t, errorStr, "username: [invalid username]")
}

func TestValidator_Clear(t *testing.T) {
	v := datacop.New()

	// Add some errors
	v.Check(false, "field1", "error1")
	v.Check(false, "field2", "error2")
	assert.True(t, v.HasErrors())

	// Clear errors
	v.Clear()
	assert.False(t, v.HasErrors())
	assert.Empty(t, v.Error())
}

func TestValidator_JSONMarshaling(t *testing.T) {
	v := datacop.New()

	// Add various types of errors
	v.Check(false, "username", "invalid username")
	v.Check(false, "password", "too short")
	v.CheckStandalone(false, "validation failed")

	// Marshal to JSON
	data, err := json.Marshal(v)
	require.NoError(t, err)

	// Verify JSON structure
	expected := map[string]map[string]string{
		"fields": {
			"username":                 "invalid username",
			"password":                 "too short",
			datacop.StandaloneErrorKey: "validation failed",
		},
	}

	var actual map[string]map[string]string
	err = json.Unmarshal(data, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestValidator_Merge(t *testing.T) {
	v1 := datacop.New()
	v2 := datacop.New()

	// Add errors to both validators
	v1.Check(false, "field1", "error1")
	v2.Check(false, "field2", "error2")
	v2.CheckStandalone(false, "standalone error")

	// Merge v2 into v1
	v1.Merge(v2)

	// Check that v1 now contains all errors
	assert.True(t, v1.HasErrorFor("field1"))
	assert.True(t, v1.HasErrorFor("field2"))
	assert.True(t, v1.HasStandaloneErrors())
}

func TestValidator_ValidationChaining(t *testing.T) {
	v := datacop.New()

	// Test chaining of validations
	tests := []struct {
		name     string
		validate func()
		expect   func(t *testing.T)
	}{
		{
			name: "single field multiple validations",
			validate: func() {
				v.Field("password", "weak123").
					Check(is.Required("weak123"), "password required").
					Check(is.MinLength(8)("weak123"), "password too short").
					Check(is.Match(`[A-Z]`)("weak123"), "needs uppercase")
			},
			expect: func(t *testing.T) {
				assert.True(t, v.HasErrorFor("password"))
				assert.Contains(t, v.ErrorFor("password"), "needs uppercase")
				assert.Contains(t, v.ErrorFor("password"), "password too short")
			},
		},
		{
			name: "multiple fields with groups",
			validate: func() {
				v.Clear()
				v.Group("user").
					Field("email", "invalid").
					Check(is.Email("invalid"), "invalid email")
			},
			expect: func(t *testing.T) {
				assert.True(t, v.HasErrorFor("user.email"))
				assert.Contains(t, v.ErrorFor("user.email"), "invalid email")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate()
			tt.expect(t)
		})
	}
}

func TestValidator_Errors(t *testing.T) {
	v := datacop.New()

	// Add various types of errors
	v.Check(false, "field1", "error1")
	v.Check(false, "field2", "error2")
	v.CheckStandalone(false, "standalone")

	// Test Errors() map
	errors := v.Errors()
	assert.Equal(t, "error1", errors["field1"])
	assert.Equal(t, "error2", errors["field2"])
	assert.Equal(t, "standalone", errors[datacop.StandaloneErrorKey])

	// Test ValidationErrors() detailed errors
	valErrors := v.ValidationErrors()
	assert.Len(t, valErrors["field1"], 1)
	assert.Equal(t, "error1", valErrors["field1"][0].Message)
}

func TestFieldValidation(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         any
		validations   func(*datacop.Validator)
		expectErrors  bool
		expectedError string
	}{
		{
			name:  "single valid field",
			field: "username",
			value: "johndoe",
			validations: func(v *datacop.Validator) {
				v.Field("username", "johndoe").
					Check(is.Required("johndoe"), "username required").
					Check(is.MinLength(3)("johndoe"), "username too short")
			},
			expectErrors: false,
		},
		{
			name:  "multiple validation failures",
			field: "password",
			value: "weak",
			validations: func(v *datacop.Validator) {
				v.Field("password", "weak").
					Check(is.MinLength(8)("weak"), "password too short").
					Check(is.Match(`[A-Z]`)("weak"), "must contain uppercase")
			},
			expectErrors:  true,
			expectedError: "password: [password too short, must contain uppercase]",
		},
		{
			name:  "valid field with multiple checks",
			field: "email",
			value: "test@example.com",
			validations: func(v *datacop.Validator) {
				v.Field("email", "test@example.com").
					Check(is.Required("test@example.com"), "email required").
					Check(is.Email("test@example.com"), "invalid email format")
			},
			expectErrors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := datacop.New()
			tt.validations(v)

			if tt.expectErrors {
				assert.True(t, v.HasErrors())
				assert.True(t, v.HasErrorFor(tt.field))
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, v.Error())
				}
			} else {
				assert.False(t, v.HasErrors())
				assert.False(t, v.HasErrorFor(tt.field))
			}
		})
	}
}

func TestGroupValidation(t *testing.T) {
	tests := []struct {
		name           string
		validations    func(*datacop.Validator)
		expectErrors   bool
		expectedErrors []string
	}{
		{
			name: "valid nested group",
			validations: func(v *datacop.Validator) {
				userGroup := v.Group("user")
				userGroup.Field("name", "John Doe").
					Check(is.Required("John Doe"), "name required")

				addressGroup := v.Group("address")
				addressGroup.Field("street", "123 Main St").
					Check(is.Required("123 Main St"), "street required")
			},
			expectErrors: false,
		},
		{
			name: "invalid nested fields",
			validations: func(v *datacop.Validator) {
				userGroup := v.Group("user")
				userGroup.Field("name", "").
					Check(is.Required(""), "name required")

				addressGroup := v.Group("address")
				addressGroup.Field("street", "").
					Check(is.Required(""), "street required")
			},
			expectErrors: true,
			//expectedError: "user.name: [name required] | address.street: [street required]",
			expectedErrors: []string{"user.name: [name required]", "address.street: [street required]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := datacop.New()
			tt.validations(v)

			if tt.expectErrors {
				assert.True(t, v.HasErrors())
				for _, err := range tt.expectedErrors {
					assert.Contains(t, v.Error(), err)
				}
			} else {
				assert.False(t, v.HasErrors())
			}
		})
	}
}

func TestWhenValidation(t *testing.T) {
	tests := []struct {
		name          string
		validations   func(*datacop.Validator)
		expectErrors  bool
		expectedError string
	}{
		{
			name: "chain: When-Check-When-Check",
			validations: func(v *datacop.Validator) {
				v.Field("role", "").
					When(true).
					Check(is.Required(""), "role required").
					When(true).
					Check(is.In("admin", "user")(""), "invalid role")
			},
			expectErrors:  true,
			expectedError: "role: [role required, invalid role]",
		},
		{
			name: "chain: When(false)-Check-When(true)-Check",
			validations: func(v *datacop.Validator) {
				v.Field("role", "").
					When(false).
					Check(is.Required(""), "role required").
					When(true).
					Check(is.In("admin", "user")(""), "invalid role")
			},
			expectErrors: false,
		},
		{
			name: "chain: When-Check-Check-When-Check",
			validations: func(v *datacop.Validator) {
				v.Field("role", "admin").
					When(true).
					Check(is.Required("admin"), "role required").
					Check(is.MinLength(3)("admin"), "role too short").
					When(true).
					Check(is.In("admin", "user")("admin"), "invalid role")
			},
			expectErrors: false,
		},
		{
			name: "chain: When-When-Check-Check",
			validations: func(v *datacop.Validator) {
				v.Field("role", "").
					When(true).
					When(false).
					Check(is.Required(""), "role required").
					Check(is.In("admin", "user")(""), "invalid role")
			},
			expectErrors: false,
		},
		{
			name: "chain: Check-When-Check (When affects only following check)",
			validations: func(v *datacop.Validator) {
				v.Field("role", "").
					Check(is.Required(""), "role required").
					When(false).
					Check(is.In("admin", "user")(""), "invalid role")
			},
			expectErrors:  true,
			expectedError: "role: [role required]",
		},
		{
			name: "chain: When-Check with mixed conditions",
			validations: func(v *datacop.Validator) {
				v.Field("role", "admin").
					When(true).
					Check(is.Required("admin"), "role required").
					When(false).
					Check(is.MinLength(10)("admin"), "role too short"). // should be skipped
					When(true).
					Check(is.In("admin", "user")("admin"), "invalid role")
			},
			expectErrors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := datacop.New()
			tt.validations(v)

			if tt.expectErrors {
				assert.True(t, v.HasErrors())
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, v.Error())
				}
			} else {
				assert.False(t, v.HasErrors())
			}
		})
	}
}
