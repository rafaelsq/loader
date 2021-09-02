package loader

func KeysToInt(slice []interface{}) []int {
	keys := make([]int, len(slice))
	for i, key := range slice {
		keys[i] = key.(int)
	}
	return keys
}

func KeysToInt64(slice []interface{}) []int64 {
	keys := make([]int64, len(slice))
	for i, key := range slice {
		keys[i] = key.(int64)
	}
	return keys
}

func KeysToString(slice []interface{}) []string {
	keys := make([]string, len(slice))
	for i, key := range slice {
		keys[i] = key.(string)
	}
	return keys
}
