package utils

func IsContain(paramet string, data []string) bool {
	for _, p := range data {
		if p == paramet {
			return true
		}
	}
	return false
}
