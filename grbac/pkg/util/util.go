package util

// Contains is used to determine if arr contains s
func Contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
