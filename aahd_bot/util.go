package aahd_bot

func FindAndDelete[T interface{}](s []T, comparator func(T) bool) []T {
	index := 0
	for _, i := range s {
		if comparator(i) {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

func Find[T interface{}](s []T, comparator func(T) bool) bool {
	for _, i := range s {
		if comparator(i) {
			return true
		}
	}
	return false
}
