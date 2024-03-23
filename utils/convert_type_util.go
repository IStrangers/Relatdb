package utils

func ConvertInt(val any) (int, bool) {
	switch v := val.(type) {
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}
