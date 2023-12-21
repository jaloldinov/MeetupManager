package slug

import (
	"strings"
)

func Make(text *string) *string {
	outputString := strings.ReplaceAll(*text, " ", "-")
	return &outputString
}
