package is

// Password returns common password validation rules.
// This is an example of a function that could be used in a project's own validation library,
// combining common validation rules into a single function.
func Password(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	return Required(str) &&
		MinLength(8)(str) &&
		Match(`[A-Z]`)(str) &&
		Match(`[a-z]`)(str) &&
		Match(`[0-9]`)(str)
}

// Username returns common username validation rules.
// This is an example of a function that could be used in a project's own validation library,
// combining common validation rules into a single function.
func Username(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	return Required(str) &&
		MinLength(3)(str) &&
		MaxLength(255)(str) &&
		Match(`^[a-zA-Z0-9_-]+$`)(str)
}
