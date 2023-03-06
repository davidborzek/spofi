package format

import "fmt"

func FormatIcon(icon string, a string) string {
	return fmt.Sprintf("%s  %s", icon, a)
}
