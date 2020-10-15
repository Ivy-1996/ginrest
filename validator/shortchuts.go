package validator

// Register your own idea by this method
// It is not goroutine safe
func RegisterLibrary(key string, validateFunc ValidateFunc) {
	simpleValidateLibrary.Register(key, validateFunc)
}
