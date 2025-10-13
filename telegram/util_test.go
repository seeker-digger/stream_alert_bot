package telegram

import "testing"

func TestChunkSlice(t *testing.T) {
	s := make([]string, 235)
	c := chunkSlice(s, 33)
	for _, i := range c {
		t.Log(len(i))
	}
}
