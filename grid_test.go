package ganim8_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yohamta/ganim8/v3"
)

func assertEqualRect(a, b *image.Rectangle) bool {
	return a.Eq(*b)
}

func assertEqualRects(a, b []*image.Rectangle) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !assertEqualRect(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestWith2Integers(t *testing.T) {
	grid := ganim8.NewGrid(16, 16, 64, 64)
	nr := func(x, y int) *image.Rectangle {
		r := image.Rect(x, y, x+16, y+16)
		return &r
	}

	var tests = []struct {
		name string
		args []interface{}
		want []*image.Rectangle
	}{
		{"returns a single frame", []interface{}{1, 1}, []*image.Rectangle{nr(0, 0)}},
		{"another single frame", []interface{}{3, 2}, []*image.Rectangle{nr(32, 16)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _TestGetFrames(grid, tt.args)
			if assertEqualRects(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestWithSeveralPairsOfIntegers(t *testing.T) {
	grid := ganim8.NewGrid(16, 16, 64, 64)
	gridWithOffesets := ganim8.NewGrid(16, 16, 64, 64, 10, 20)

	nr := func(x, y int) *image.Rectangle {
		r := image.Rect(x, y, x+16, y+16)
		return &r
	}

	var tests = []struct {
		name string
		args []interface{}
		want []*image.Rectangle
		grid *ganim8.Grid
	}{
		{"returns a list of frames", []interface{}{1, 3, 2, 2, 3, 1}, []*image.Rectangle{nr(0, 32), nr(16, 16), nr(32, 0)}, grid},
		{"takes into account left and top", []interface{}{1, 3, 2, 2, 3, 1}, []*image.Rectangle{nr(10, 52), nr(26, 36), nr(42, 20)}, gridWithOffesets},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _TestGetFrames(tt.grid, tt.args)
			if assertEqualRects(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestWithStringAndAIntegry(t *testing.T) {
	grid := ganim8.NewGrid(16, 16, 64, 64)

	nr := func(x, y int) *image.Rectangle {
		r := image.Rect(x, y, x+16, y+16)
		return &r
	}

	var tests = []struct {
		name string
		args []interface{}
		want []*image.Rectangle
		grid *ganim8.Grid
	}{
		{"returns a list of frames", []interface{}{"1-2", 2}, []*image.Rectangle{nr(0, 16), nr(16, 16)}, grid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _TestGetFrames(tt.grid, tt.args)
			if assertEqualRects(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestWithSeveralStrings(t *testing.T) {
	grid := ganim8.NewGrid(16, 16, 64, 64)

	nr := func(x, y int) *image.Rectangle {
		r := image.Rect(x, y, x+16, y+16)
		return &r
	}

	var tests = []struct {
		name string
		args []interface{}
		want []*image.Rectangle
		grid *ganim8.Grid
	}{
		{"returns a list of frames", []interface{}{"1-2", 2, 3, 2}, []*image.Rectangle{nr(0, 16), nr(16, 16), nr(32, 16)}, grid},
		{"parses rows first, then columns", []interface{}{"1-3", "1-3"}, []*image.Rectangle{
			nr(0, 0), nr(16, 0), nr(32, 0),
			nr(0, 16), nr(16, 16), nr(32, 16),
			nr(0, 32), nr(16, 32), nr(32, 32),
		}, grid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _TestGetFrames(tt.grid, tt.args)
			if assertEqualRects(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestWithoutArgs(t *testing.T) {
	grid := ganim8.NewGrid(16, 16, 64, 64)
	nr := func(x, y int) *image.Rectangle {
		r := image.Rect(x, y, x+16, y+16)
		return &r
	}

	var tests = []struct {
		name string
		args []interface{}
		want []*image.Rectangle
	}{
		{"returns all frames", []interface{}{}, []*image.Rectangle{
			nr(0, 0), nr(16, 0), nr(32, 0), nr(48, 0),
			nr(0, 16), nr(16, 16), nr(32, 16), nr(48, 16),
			nr(0, 32), nr(16, 32), nr(32, 32), nr(48, 32),
			nr(0, 48), nr(16, 48), nr(32, 48), nr(48, 48),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _TestGetFrames(grid, tt.args)
			if assertEqualRects(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestFramesWithBorder(t *testing.T) {
	gs := ganim8.NewGrid(32, 32, 1024, 1025, 0, 0, 1)
	f := gs.Frames("1-3", 1, "3-1", 1)
	require.Equal(t, 6, len(f))

	expectedFrames := []image.Rectangle{
		image.Rect(1, 1, 33, 33),
		image.Rect(34, 1, 66, 33),
		image.Rect(67, 1, 99, 33),
	}

	for i, want := range []image.Rectangle{
		expectedFrames[0],
		expectedFrames[1],
		expectedFrames[2],
		expectedFrames[2],
		expectedFrames[1],
		expectedFrames[0],
	} {
		t.Run(fmt.Sprintf("frame-%d", i), func(t *testing.T) {
			require.Equal(t, want, f[i].Bounds())
		})
	}
}

func _TestGetFrames(g *ganim8.Grid, args []interface{}) []*image.Rectangle {
	return g.GetFrames(args...)
}
