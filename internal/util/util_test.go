package util

import (
	"fmt"
	l2 "main.go/internal/logger"
	"strings"
	"testing"
)

func TestGetNLastLines(t *testing.T) {
	l, err := GetNLastLines(l2.LogFile, 3)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(strings.Join(l, "\n"))
}
