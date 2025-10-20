package util

import (
	"fmt"
	"strings"
	"testing"

	l2 "main.go/internal/logger"
)

func TestGetNLastLines(t *testing.T) {
	l, err := GetNLastLines(l2.LogFile, 3)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(strings.Join(l, "\n"))
}
