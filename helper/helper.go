package helper

const MaxQuadSize = 20

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c rune) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}

func IsBlank(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
