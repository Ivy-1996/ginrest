package validator

func RegisterLibrary(key string, validateFunc ValidateFunc) {
	simpleValidateLibrary.Register(key, validateFunc)
}
