package helper

func IsEmpty(param string) bool {
	return param == ""
}

func ErrorIsNil(err error) bool {
	return err != nil
}

func NotMatch(a, b string) bool {
	return a != b
}
