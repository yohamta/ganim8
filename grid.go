package ganim8

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"regexp"
	"strconv"
)

func assertPositiveInteger(value int, name string) {
	if value < 1 {
		log.Fatal(fmt.Sprintf("%s should be a positive number, was %d", name, value))
	}
}

func assertSize(size, limit int, name string) {
	if size > limit {
		log.Fatal(fmt.Sprintf("%s should be <= %d, was %d", name, limit, size))
	}
}

type frameCache map[string]map[int]map[int]*image.Rectangle

var _frames frameCache

var intervalMatcher regexp.Regexp

func init() {
	_frames = make(map[string]map[int]map[int]*image.Rectangle)
	intervalMatcher = *regexp.MustCompile("^([0-9]+)-([0-9]+)$")
}

// Grid represents a grid
type Grid struct {
	frameWidth, frameHeight int
	imageWidth, imageHeight int
	left, top               int
	width, height           int
	border                  int
	key                     string
}

// NewGrid returns a new grid with specified frame size, image size, and
// offsets (left, top).
//
// Grids have only one purpose: To build groups of quads of the same size as
// easily as possible.
//
// They need to know only 2 things: the size of each frame and the size of
// the image they will be applied to.
//
// Grids are just a convenient way of getting frames from a sprite.
// Frames are assumed to be distributed in rows and columns.
// Frame 1,1 is the one in the first row, first column.
func NewGrid(frameWidth, frameHeight, imageWidth, imageHeight int, args ...int) *Grid {
	assertPositiveInteger(frameWidth, "frameWidth")
	assertPositiveInteger(frameHeight, "frameHeight")
	assertPositiveInteger(imageWidth, "imageWidth")
	assertPositiveInteger(imageHeight, "imageHeight")
	assertSize(frameWidth, imageWidth, "frameWidth")
	assertSize(frameHeight, imageHeight, "frameHeight")

	left, top, border := 0, 0, 0
	switch len(args) {
	case 3:
		border = args[2]
		fallthrough
	case 2:
		top = args[1]
		fallthrough
	case 1:
		left = args[0]
	}

	g := &Grid{
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		imageWidth:  imageWidth,
		imageHeight: imageHeight,
		left:        left,
		top:         top,
		width:       imageWidth / frameWidth,
		height:      imageHeight / frameHeight,
		border:      border,
	}

	g.key = getGridKey(g.frameWidth, g.frameHeight, g.imageWidth,
		g.imageHeight, g.left, g.top)

	return g
}

func getGridKey(args ...int) string {
	var b bytes.Buffer
	s := ""
	for _, a := range args {
		b.Write([]byte(s))
		b.Write([]byte(strconv.Itoa(a)))
		s = "-"
	}
	return b.String()
}

func (g *Grid) createFrame(x, y int) *image.Rectangle {
	fw, fh := g.frameWidth, g.frameHeight
	x0, y0 := g.left+(x-1)*fw+x*g.border, g.top+(y-1)*fh+y*g.border
	r := image.Rect(x0, y0, x0+fw, y0+fh)
	return &r
}

func (g *Grid) getOrCreateFrame(x, y int) *image.Rectangle {
	if x < 1 || x > g.width || y < 1 || y > g.height {
		log.Fatal(fmt.Sprintf("There is no frame for x=%d, y=%d", x, y))
	}
	key := g.key
	if _, ok := _frames[key]; !ok {
		_frames[key] = map[int]map[int]*image.Rectangle{}
	}
	if _, ok := _frames[key][x]; !ok {
		_frames[key][x] = map[int]*image.Rectangle{}
	}
	if _, ok := _frames[key][x][y]; !ok {
		_frames[key][x][y] = g.createFrame(x, y)
	}
	return _frames[key][x][y]
}

// GetFrames accepts an arbitrary number of parameters.
// They can be either numbers or strings.
//
// Each two numbers are interpreted as quad coordinates in the
// format (column, row). This way, grid:getFrames(3,4) will return
// the frame in column 3, row 4 of the grid.
//
// There can be more than just two: grid:getFrames(1,1, 1,2, 1,3)
// will return the frames in {1,1}, {1,2} and {1,3} respectively.
func (g *Grid) GetFrames(args ...interface{}) []*image.Rectangle {
	result := []*image.Rectangle{}
	if len(args) == 0 {
		for y := 1; y <= g.height; y++ {
			for x := 1; x <= g.width; x++ {
				result = append(result, g.getOrCreateFrame(x, y))
			}
		}
		return result
	}
	for i := 0; i < len(args); i += 2 {
		minx, maxx, stepx := parseInterval(args[i])
		miny, maxy, stepy := parseInterval(args[i+1])
		for y := miny; stepy > 0 && y <= maxy || stepy < 0 && y >= maxy; y += stepy {
			for x := minx; stepx > 0 && x <= maxx || stepx < 0 && x >= maxx; x += stepx {
				result = append(result, g.getOrCreateFrame(x, y))
			}
		}
	}
	return result
}

// Width returns the width of the grid
func (g *Grid) Width() int {
	return g.width
}

// Height returns the height of the grid
func (g *Grid) Height() int {
	return g.height
}

// Frames is a shorter name of GetFrames
func (g *Grid) Frames(args ...interface{}) []*image.Rectangle {
	return g.GetFrames(args...)
}

// G is a shorter name of GetFrames
func (g *Grid) G(args ...interface{}) []*image.Rectangle {
	return g.GetFrames(args...)
}
