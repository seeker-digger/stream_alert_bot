package util

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// GetNLastLines returns the last n lines of the file located at filePath.
//
// It reads the file sequentially using bufio.Scanner and maintains a rolling
// buffer of up to n lines, so memory usage is proportional to n (not the
// entire file size). The returned slice contains the lines in the same order
// they appear in the file (oldest of the returned lines first, newest last).
//
// If the file contains fewer than n lines, all available lines are returned.
// Lines are returned without trailing newline characters (bufio.Scanner.Text()).
//
// Errors returned wrap underlying I/O errors encountered when opening or
// scanning the file. The parameter n must be positive; passing 0 or a
// negative value will cause a runtime panic due to an invalid slice capacity.
// Note: bufio.Scanner has a default maximum token size (around 64KB); very
// long lines may require a different reading strategy.
func GetNLastLines(filePath string, n int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error in opening file: %w", err)
	}
	defer file.Close()

	lines := make([]string, 0, n)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if len(lines) < n {
			lines = append(lines, line)
		} else {
			copy(lines, lines[1:])

			lines[n-1] = line
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error in scanning file: %w", err)
	}

	return lines, nil
}

// WaitForSignal just a stub for now.
func WaitForSignal() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	os.Exit(1)

}
