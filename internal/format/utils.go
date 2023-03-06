package format

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func substring(s string, from int, to int) string {
	return string([]rune(s)[from:to])
}

func getMaxLengths(data [][]string, maximum int) []int {
	arrSize := 0
	for _, d := range data {
		if len(d) > arrSize {
			arrSize = len(d)
		}
	}

	maxLengths := make([]int, arrSize)
	for _, d := range data {
		for i, t := range d {
			l := utf8.RuneCountInString(t)

			if l > maxLengths[i] {
				maxLengths[i] = l
			}

			if maxLengths[i] > maximum {
				maxLengths[i] = maximum
			}
		}
	}
	return maxLengths
}

func buildRow(data []string, maxLengths []int, lengthLimit int) string {
	out := ""
	for i, t := range data {
		l := utf8.RuneCountInString(t)
		if l > lengthLimit-3 {
			l = lengthLimit
			t = fmt.Sprintf("%s...", substring(t, 0, lengthLimit-3))
		}

		spacer := ""
		if i != len(data)-1 {
			diff := maxLengths[i] - l
			spacer = strings.Repeat(" ", diff+10)
		}

		out = fmt.Sprintf("%s%s%s",
			out,
			t,
			spacer,
		)
	}

	return out
}

func BuildRows(data [][]string, maxColumnSize int) []string {
	maxLengths := getMaxLengths(data, maxColumnSize)

	out := make([]string, len(data))
	for i, d := range data {
		out[i] = buildRow(d, maxLengths, maxColumnSize)
	}
	return out
}
