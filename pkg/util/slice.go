package util

func SliceEqual[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	if (s1 == nil) != (s2 == nil) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// IsSubSlice returns whether src is the sub-slice of dst
func IsSubSlice[T comparable](dst, src []T) bool {
	if len(src) > len(dst) {
		return false
	}

	if (src == nil) != (dst == nil) {
		return false
	}

	m := make(map[T]struct{}, len(dst))
	for _, v := range dst {
		m[v] = struct{}{}
	}

	for _, v := range src {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}
