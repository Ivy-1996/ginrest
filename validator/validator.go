package validator

type ValidFunc func(value interface{}) bool

func Validate(value string, validFunc ValidFunc) bool {
	return validFunc(value)
}

func InterRequired(value string) bool {
	return InterRegx.MatchString(value)
}

func EmailRequired(value string) bool {
	return EmailRegex.MatchString(value)
}

func UuidRequired(value string) bool {
	return UUIDRegex.MatchString(value)
}
