package telegram

import (
	"reflect"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	s := make([]string, 235)
	c := chunkSlice(s, 33)
	for _, i := range c {
		t.Log(len(i))
	}
}

func TestRemoveAll_Ints(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		value int
		want  []int
	}{
		{"nil slice", nil, 1, []int(nil)},
		{"empty slice", []int{}, 1, []int{}},
		{"no match", []int{1, 2, 3}, 4, []int{1, 2, 3}},
		{"some matches", []int{1, 2, 1, 3, 1}, 1, []int{2, 3}},
		{"all match", []int{5, 5, 5}, 5, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeAll(tt.slice, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("removeAll(%v, %v) = %v; want %v", tt.slice, tt.value, got, tt.want)
			}
		})
	}
}

func TestRemoveAll_Strings(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		value string
		want  []string
	}{
		{"nil slice", nil, "a", []string(nil)},
		{"empty slice", []string{}, "a", []string{}},
		{"no match", []string{"a", "b"}, "c", []string{"a", "b"}},
		{"some matches", []string{"a", "b", "a", "c"}, "a", []string{"b", "c"}},
		{"all match", []string{"x", "x"}, "x", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeAll(tt.slice, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("removeAll(%v, %q) = %v; want %v", tt.slice, tt.value, got, tt.want)
			}
		})
	}
}

func TestRemoveAll_ComparableStructs(t *testing.T) {
	type P struct{ A, B int }
	a := P{1, 2}
	b := P{3, 4}
	tests := []struct {
		name  string
		slice []P
		value P
		want  []P
	}{
		{"mixed structs", []P{a, b, a}, a, []P{b}},
		{"none match", []P{a, b}, P{9, 9}, []P{a, b}},
		{"all match", []P{b, b}, b, []P{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeAll(tt.slice, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("removeAll(%v, %v) = %v; want %v", tt.slice, tt.value, got, tt.want)
			}
		})
	}
}
