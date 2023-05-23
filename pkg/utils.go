package pkg

func IsTheSameArray(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	var tmp = make(map[string]string, len(a))
	for _, el := range a {
		tmp[el] = el
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
}
