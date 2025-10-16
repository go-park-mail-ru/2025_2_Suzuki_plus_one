package utils

// Function that helps safely logging token
func SafeTokenPrefix(tokenString string) string {
	if (len(tokenString)) < 8 {
		return "[short]"
	}

	return tokenString[:8] + "..."
}
