package util

import (
	"strings"
)

func TextToLines(t string, width int) []string {
	var lines []string
	nonEmptyLine := false // flag to track whether the previous line had content

	for _, line := range strings.Split(t, "\n") {
		// Split the line into multiple lines if it is too long
		for len(line) > width {
			// Find the index of the last space character before the max width
			spaceIdx := strings.LastIndex(line[:width+1], " ")

			// If no space character was found, break the word at the max width
			if spaceIdx == -1 {
				spaceIdx = width
			}

			// Add the line to the output and remove it from the input
			lines = append(lines, line[:spaceIdx])
			line = line[spaceIdx+1:]
			nonEmptyLine = true
		}

		// Add the line to the output if it is not empty
		if len(line) > 0 {
			lines = append(lines, line)
			nonEmptyLine = true
		}

		// Add an empty line to the output if the previous line had content
		if nonEmptyLine {
			lines = append(lines, "")
			nonEmptyLine = false
		}
	}

	// Remove the trailing empty line from the output
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
