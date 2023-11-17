package hello

import (
	"strings"
)

func Say(str []string) string {
	if len(str) == 0 {
		str = []string{"World"}
	}

	return "Hello " + strings.Join(str, ", ") + "!"
}
