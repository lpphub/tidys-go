package slices

// collection 通用函数

// Map 泛型转换
func Map[T any, R any](arr []T, fn func(T) R) []R {
	if arr == nil {
		return []R{}
	}
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

// Filter 泛型过滤
func Filter[T any](arr []T, fn func(T) bool) []T {
	if arr == nil {
		return []T{}
	}
	result := make([]T, 0, len(arr))
	for _, v := range arr {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// GroupBy 将切片按 key 函数分组
func GroupBy[K comparable, V any](items []V, getKey func(V) K) map[K][]V {
	grouped := make(map[K][]V)
	for _, item := range items {
		key := getKey(item)
		grouped[key] = append(grouped[key], item)
	}
	return grouped
}

// IndexBy 将切片按 key 构建为 Map
func IndexBy[K comparable, V any](items []V, getKey func(V) K) map[K]V {
	m := make(map[K]V, len(items))
	for i := range items {
		m[getKey(items[i])] = items[i]
	}
	return m
}
