package easypost

// StringPtr returns a pointer to a string with the given value.
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to a bool with the given value.
func BoolPtr(b bool) *bool {
	return &b
}
