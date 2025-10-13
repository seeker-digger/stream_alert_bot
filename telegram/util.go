package telegram

import "strings"

func removeAll[T comparable](slice []T, value T) []T {
	result := slice[:0]
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

func chunkSlice[T any](slice []T, size int) [][]T {
	if size <= 0 {
		panic("size must be positive")
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func scheduleUpd() {
	go func() {
		needUpdate <- true
	}()
}

func escapeMarkdownV2Text(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		"!", "\\!",
	)
	return replacer.Replace(s)
}
