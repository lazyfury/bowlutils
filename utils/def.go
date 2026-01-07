package utils

func Def[T comparable](val T, defaultValue T) T {
	if IsZero(val) {
		return defaultValue
	}
	return val
}
