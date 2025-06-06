package utils

func Contains[T comparable](s []T, i T) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}

func Deduplicate[T comparable](list []T) []T {
	seen := make(map[T]bool)
	var filtered []T
	for _, val := range list {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			filtered = append(filtered, val)
		}
	}
	return filtered
}

// 求并集
func Union[T comparable](slice1, slice2 []T) []T {
	m := make(map[T]int)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// 求交集
func Intersect[T comparable](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func Difference[T comparable](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
