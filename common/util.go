package common

// HasString check if the array has input string item
func HasString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

// GetFirstStringInMap get first string value in string map
func GetFirstStringInMap(input map[string]string) string {
	for _, v := range input {
		return v
	}
	return ""
}
