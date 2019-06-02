package common

// ItemInList util
func ItemInList(a string, b []string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}
