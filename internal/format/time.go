package format

import (
	"fmt"
	"strconv"
	"time"
)

// FormatTime formats a timestamp in ms
// to mm:ss format with leading zero for seconds.
func FormatTime(ms int) string {
	t := time.UnixMilli(int64(ms))
	minute := strconv.Itoa(t.Minute())

	second := strconv.Itoa(t.Second())
	if t.Second() < 10 {
		second = fmt.Sprintf("0%d", t.Second())
	}

	return fmt.Sprintf("%s:%s", minute, second)
}
