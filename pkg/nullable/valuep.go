package nullable

func StringP(s string) *string {
	return &s
}

func Int32P(i int32) *int32 {
	return &i
}

func BoolP(b bool) *bool {
	return &b
}
