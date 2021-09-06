package loader

func KeysToInt(slice []interface{}) []int {
	check := make(map[int]struct{}, len(slice))

	keys := make([]int, 0, len(slice))
	for _, key := range slice {
		value := key.(int)
		if _, exists := check[value]; !exists {
			keys = append(keys, value)
			check[value] = struct{}{}
		}
	}

	return keys
}

func KeysToInt64(slice []interface{}) []int64 {
	check := make(map[int64]struct{}, len(slice))

	keys := make([]int64, 0, len(slice))
	for _, key := range slice {
		value := key.(int64)
		if _, exists := check[value]; !exists {
			keys = append(keys, value)
			check[value] = struct{}{}
		}
	}

	return keys
}

func KeysToString(slice []interface{}) []string {
	check := make(map[string]struct{}, len(slice))

	keys := make([]string, 0, len(slice))
	for _, key := range slice {
		value := key.(string)
		if _, exists := check[value]; !exists {
			keys = append(keys, value)
			check[value] = struct{}{}
		}
	}

	return keys
}
