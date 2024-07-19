package format

import "fmt"

type Keybinding struct {
	Key         string
	Description string
}

func FormatKeybindings(
	keys ...Keybinding,
) string {
	var str string
	for idx, key := range keys {
		str += fmt.Sprintf("<b>%s:</b> %s", key.Key, key.Description)
		if idx != len(keys)-1 {
			str += " | "
		}
	}

	return str
}
