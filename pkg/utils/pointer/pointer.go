package pointer

// reference: https://github.com/kubernetes/utils

// Int32Ptr returns a pointer to an int64
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr returns a pointer to an bool
func BoolPtr(b bool) *bool {
	return &b
}
