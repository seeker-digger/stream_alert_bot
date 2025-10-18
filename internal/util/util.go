package util

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

func WaitForSignal() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	_ = <-quit
	os.Exit(1)

}
