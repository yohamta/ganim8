package ganim8

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type imageCache map[*ebiten.Image]map[*image.Rectangle]*ebiten.Image

var _imageCache imageCache

func init() {
	_imageCache = make(map[*ebiten.Image]map[*image.Rectangle]*ebiten.Image)
}

func parseDurations(durations interface{}, frameCount int) []time.Duration {
	result := make([]time.Duration, frameCount)
	switch val := durations.(type) {
	case time.Duration:
		for i := 0; i < frameCount; i++ {
			result[i] = val
		}
	case []time.Duration:
		for i := range val {
			result[i] = val[i]
		}
	case map[string]time.Duration:
		for key, duration := range val {
			min, max, step := parseInterval(key)
			for i := min; i <= max; i += step {
				result[i-1] = duration
			}
		}
	default:
		log.Fatal(fmt.Sprintf("durations must be time.Duration or []time.Duration or map[string]time.Duration. was %v", durations))
	}
	return result
}

func parseIntervals(durations []time.Duration) ([]time.Duration, time.Duration) {
	result := []time.Duration{0}
	var time time.Duration = 0
	for _, v := range durations {
		time += v
		result = append(result, time)
	}
	return result, time
}

// Status represents the animation status.
type Status int

const (
	Playing = iota
	Paused
)

// Animation represents an animation created from specified frames
// and an *ebiten.Image
type Animation struct {
	sprite             *Sprite
	position           int
	timer              time.Duration
	durations          []time.Duration
	intervals          []time.Duration
	totalDuration      time.Duration
	onLoop             OnLoop
	status             Status
	flippedH, flippedV bool
}

// OnLoop is callback function which representing
// one of the animation methods.
// it will be called every time an animation "loops".
//
// It will have two parameters: the animation instance,
// and how many loops have been elapsed.
//
// The value would be Nop (No operation) if there's nothing
// to do except for looping the animation.
//
// The most usual value (apart from none) is the string 'pauseAtEnd'.
// It will make the animation loop once and then pause
// and stop on the last frame.
type OnLoop func(anim *Animation, loops int)

// Nop does nothing.
func Nop(anim *Animation, loops int) {}

// Pause pauses the animation on loop finished.
func Pause(anim *Animation, loops int) {
	anim.Pause()
}

// PauseAtEnd pauses the animation and set the position to
// the last frame.
func PauseAtEnd(anim *Animation, loops int) {
	anim.PauseAtEnd()
}

// PauseAtStart pauses the animation and set the position to
// the first frame.
func PauseAtStart(anim *Animation, loops int) {
	anim.PauseAtStart()
}

// NewAnimation returns a new animation object
//
// durations is a time.Duration or a []time.Duration or
// a map[string]time.Duration.
// When it's a time.Duration, it represents the duration of
// all frames in the animation.
// When it's a []time.Duration, it can represent different
// durations for different frames.
// You can specify durations for all frames individually,
// like this: []time.Duration { 100 * time.Millisecond,
// 100 * time.Millisecond } or you can specify durations for
// ranges of frames: map[string]time.Duration { "1-2":
// 100 * time.Millisecond, "3-5": 200 * time.Millisecond }.
func NewAnimation(sprite *Sprite, durations interface{}, onLoop OnLoop) *Animation {
	_durations := parseDurations(durations, sprite.length)
	intervals, totalDuration := parseIntervals(_durations)
	anim := &Animation{
		sprite:        sprite,
		position:      0,
		timer:         0,
		durations:     _durations,
		intervals:     intervals,
		totalDuration: totalDuration,
		onLoop:        onLoop,
		status:        Playing,
		flippedH:      false,
		flippedV:      false,
	}
	return anim
}

// Clone return a copied animation object.
func (anim *Animation) Clone() *Animation {
	new := *anim
	return &new
}

// FlipH flips the animation horizontally.
func (anim *Animation) FlipH() {
	anim.flippedH = !anim.flippedH
}

// FlipV flips the animation horizontally.
func (anim *Animation) FlipV() {
	anim.flippedV = !anim.flippedV
}

func seekFrameIndex(intervals []time.Duration, timer time.Duration) int {
	high, low, i := len(intervals)-2, 0, 0
	for low <= high {
		i = (low + high) / 2
		if timer >= intervals[i+1] {
			low = i + 1
		} else if timer < intervals[i] {
			high = i - 1
		} else {
			return i
		}
	}
	return i
}

// Update updates the animation.
func (anim *Animation) Update(elapsedTime time.Duration) {
	if anim.status != Playing || anim.sprite.length <= 1 {
		return
	}
	anim.timer += elapsedTime
	loops := anim.timer / anim.totalDuration
	if loops != 0 {
		anim.timer = anim.timer - anim.totalDuration*loops
		(anim.onLoop)(anim, int(loops))
	}
	anim.position = seekFrameIndex(anim.intervals, anim.timer)
}

// Status returns the status of the animation.
func (anim *Animation) Status() Status {
	return anim.status
}

// Pause pauses the animation.
func (anim *Animation) Pause() {
	anim.status = Paused
}

// Position returns the current position of the frame.
// The position counts from 1 (not 0).
func (anim *Animation) Position() int {
	return anim.position + 1
}

// Duration returns the current durations of each frames.
func (anim *Animation) Durations() []time.Duration {
	return anim.durations
}

// TotalDuration returns the total duration of the animation.
func (anim *Animation) TotalDuration() time.Duration {
	return anim.totalDuration
}

// Size returns the size of the current frame.
func (anim *Animation) Size() (float64, float64) {
	return anim.sprite.Size()
}

// W is a shortcut for Size().X.
func (anim *Animation) W() float64 {
	return anim.sprite.W()
}

// H is a shortcut for Size().Y.
func (anim *Animation) H() float64 {
	return anim.sprite.H()
}

// Timer returns the current accumulated times of current frame.
func (anim *Animation) Timer() time.Duration {
	return anim.timer
}

// GoToFrame sets the position of the animation and
// sets the timer at the start of the frame.
func (anim *Animation) GoToFrame(position int) {
	anim.position = position - 1
	anim.timer = anim.intervals[anim.position]
}

// PauseAtEnd pauses the animation and set the position
// to the last frame.
func (anim *Animation) PauseAtEnd() {
	anim.position = anim.sprite.length - 1
	anim.timer = anim.totalDuration
	anim.Pause()
}

// PauseAtStart pauses the animation and set the position
// to the first frame.
func (anim *Animation) PauseAtStart() {
	anim.position = 0
	anim.timer = 0
	anim.status = Paused
}

// Resume resumes the animation
func (anim *Animation) Resume() {
	anim.status = Playing
}

// Draw draws the animation with the specified option parameters.
func (anim *Animation) Draw(screen *ebiten.Image, opts *DrawOptions) {
	anim.sprite.Draw(screen, anim.position, opts)
}

// DrawWithShader draws the animation with the specified option parameters.
func (anim *Animation) DrawWithShader(screen *ebiten.Image, opts *DrawOptions, shaderOpts *ShaderOptions) {
	anim.sprite.DrawWithShader(screen, anim.position, opts, shaderOpts)
}
