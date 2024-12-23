package is

import "regexp"

const (
	rgxEmail = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	rgxPhone = `^\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})$`
)

// Email is a very simple email validation function. For a more comprehensive
// email validation, consider using a package like github.com/patrickward/mailcop.
//
// Example usage:
// Email("foo@example.com") // returns true
// Email("invalid-email") // returns false
func Email(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return regexp.MustCompile(rgxEmail).MatchString(str)
}

// Phone is a simple phone number validation function. It expects a string
// with a format of 123-456-7890.
func Phone(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return regexp.MustCompile(rgxPhone).MatchString(str)
}
