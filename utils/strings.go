package utils

func StringPtr(s string) *string {
	return &s
}

func StringValue(s *string) string {
	return *s
}
