package chronicle

import (
	"os"
)

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func envSearch(s string) string {
	if v := os.Getenv(s); v != "" {
		return v
	}

	return ""
}

// checks if a string is present in a slice.
func contains(s []string, s1 string) bool {
	for _, a := range s {
		if a == s1 {
			return true
		}
	}
	return false
}
