package ganim8_test

import (
	"image"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8"
)

var mockImg = ebiten.NewImage(1, 1)

func mockFrames(n int) []*image.Rectangle {
	result := make([]*image.Rectangle, n)
	for i := 0; i < n; i++ {
		r := image.Rect(i, 0, 16+i*16, 16)
		result[i] = &r
	}
	return result
}

func mockSprite(n int) *ganim8.Sprite {
	f := mockFrames(n)
	spr := ganim8.NewSprite(mockImg, f)
	return spr
}

func assertEqualDurations(a, b []time.Duration) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParsingDuration(t *testing.T) {
	var tests = []struct {
		name string
		args time.Duration
		want []time.Duration
	}{
		{"reads a simple array", 3, []time.Duration{3, 3, 3, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), tt.args, ganim8.Nop)
			got := anim.Durations()
			if assertEqualDurations(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParsingDurationArr(t *testing.T) {
	var tests = []struct {
		name string
		args []time.Duration
		want []time.Duration
	}{
		{"reads a simple val", []time.Duration{
			1, 2, 3, 4,
		}, []time.Duration{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), tt.args, ganim8.Nop)
			got := anim.Durations()
			if assertEqualDurations(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParsingDurationHash(t *testing.T) {
	var tests = []struct {
		name string
		args map[string]time.Duration
		want []time.Duration
	}{
		{"reads a simple hash", map[string]time.Duration{
			"1": 1, "2-3": 2, "4": 3,
		}, []time.Duration{1, 2, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), tt.args, ganim8.Nop)
			got := anim.Durations()
			if assertEqualDurations(got, tt.want) == false {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTotalDuration(t *testing.T) {
	var tests = []struct {
		name string
		args []time.Duration
		want time.Duration
	}{
		{"sums up the total duration", []time.Duration{1, 2, 3, 4}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), tt.args, ganim8.Nop)
			got := anim.TotalDuration()
			if got != tt.want {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}
	var tests = []struct {
		name           string
		durations      interface{}
		updateDuration time.Duration
		want           int
	}{
		{"moves to the next frame", d(1), d(1), 2},
		{"moves several mockSprite if needed", d(1), d(2), 3},
		{"when the last frame is spent goes back to the first frame",
			d(1), d(4), 1},
		{"when there're different durations per frame",
			[]time.Duration{d(1), d(2), d(3), d(4)}, d(3), 3},
		{"when there're different durations per frame",
			[]time.Duration{d(1), d(2), d(3), d(4)}, d(6), 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), tt.durations, ganim8.Nop)
			anim.Update(tt.updateDuration)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestCallback(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}

	var tests = []struct {
		name string
		arg  time.Duration
		want int
	}{
		{"invokes the onloop callback", d(4), 1},
		{"counts the loops", d(8), 2},
		{"counts negative loops", -d(4), -1},
		{"counts negative loops", -d(8), -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := 0
			anim := ganim8.NewAnimation(mockSprite(4), d(1),
				func(anim *ganim8.Animation, loops int) {
					got += loops
				})
			anim.Update(tt.arg)
			if got != tt.want {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPause(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}

	var tests = []struct {
		name string
		arg  time.Duration
		want int
	}{
		{"stops animations from happening", d(2), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), d(1), ganim8.Nop)
			anim.Pause()
			anim.Update(tt.arg)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestResume(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}
	var tests = []struct {
		name string
		arg  time.Duration
		want int
	}{
		{"resume paused animations", d(2), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), d(1), ganim8.Nop)
			anim.Pause()
			anim.Resume()
			anim.Update(tt.arg)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("%s: got %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGotoFrame(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}
	var tests = []struct {
		name  string
		arg   int
		want  int
		want2 time.Duration
	}{
		{"moves to the position and time to the frame specified", 3, 3, d(2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), d(1), ganim8.Nop)
			anim.GoToFrame(tt.arg)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("position - %s: got %v; want %v", tt.name, got, tt.want)
			}
			got2 := anim.Timer()
			if got2 != tt.want2 {
				t.Errorf("timer - %s: got %v; want %v", tt.name, got2, tt.want2)
			}
		})
	}
}

func TestPauseAtEnd(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}
	var tests = []struct {
		name  string
		arg   time.Duration
		want  int
		want2 ganim8.Status
	}{
		{"goes to the last frame and pause", d(4), 4, ganim8.Paused},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), d(1), ganim8.PauseAtEnd)
			anim.Update(tt.arg)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("positon - %s: got %v; want %v", tt.name, got, tt.want)
			}
			got2 := anim.Status()
			if got2 != tt.want2 {
				t.Errorf("status - %s: got %v; want %v", tt.name, got2, tt.want2)
			}
		})
	}
}

func TestPauseAtStart(t *testing.T) {
	d := func(s time.Duration) time.Duration {
		return s * time.Second
	}
	var tests = []struct {
		name  string
		arg   time.Duration
		want  int
		want2 ganim8.Status
	}{
		{"goes to the first frame and pause", d(4), 1, ganim8.Paused},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anim := ganim8.NewAnimation(mockSprite(4), d(1), ganim8.PauseAtStart)
			anim.Update(tt.arg)
			got := anim.Position()
			if got != tt.want {
				t.Errorf("positon - %s: got %v; want %v", tt.name, got, tt.want)
			}
			got2 := anim.Status()
			if got2 != tt.want2 {
				t.Errorf("status - %s: got %v; want %v", tt.name, got2, tt.want2)
			}
		})
	}
}
