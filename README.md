# ganim8

ganim8 is an Animation library for [Ebiten](https://ebiten.org/) which is a golang version of [anim8](https://github.com/kikito/anim8).

[GoDoc](https://pkg.go.dev/github.com/yohamta/ganim8)

## Example

```go
type Game struct {
	prevUpdateTime time.Time
	anim           *ganim8.Animation
}

func NewGame() *Game {
	g := &Game{ prevUpdateTime: time.Now() }

	// Grids are just a convenient way of getting frames of the same size from a
	// single sprite texture.
	// animation frames are represented as groups of rectangles (image.Rectangle).
	// They need to know only 2 things: the size of frames and the size of
	// the image they will be applied to, and those are the first 4 parameters of NewGrid()
	// 
	// NewGrid() parameters accpets 4 - 6 parameters:
	// ganim8.NewGrid(frameWidth, frameHeight, imageWidth, imageHeight, left, top)
	// "left" and "top" are optional, and both default to 0. They are "the left
	// and top coordinates of the point in the image where you want to put the
	// origin of coordinates of the grid".
	//
	// In this example, it make a grid with the frame size of {100,100} and
	// the image size of {500,600}
	grid := ganim8.NewGrid(100, 100, 500, 600)

	// Grids have one important method: Grid.GetFrames(...).
	// Grid.GetFrames() accepts an arbitrary number of parameters.
	// They can be either number or strings.
	// Each two numbers are interpreted as rectangle coordinates in
	// the format (column, row).
	// 
	// This way, grid.GetFrames(3, 4) will return the frame in column 3,
	// row 4 of the grid.
	// They can be more than just two: grid.GetFrames(1,1, 1,2, 1,3) will
	// return frames in {1,2}, {1,2} and {1,3} respectively.
	//
	// Using numbers for long rows or columns is tedious - so grid
	// also accpet strings indicating range plus a row/column index.
	// A row can be fetch by calling grid.GetFrames('range', rowNumber) and
	// a column by calling grid.GetFrames(columnNumber, 'range').
	//
	// It's also possible to combine both formats.
	// For example: grid.GetFrames(1, 3, 1, '1-3') will get the frame in {1,4}
	// plus the frames 1 to 3 in column 1.
	// 
	// The below code get frames 1 to 5 column in row 5
	frames := grid.GetFrames("1-5", 5)

	// Animations are groups of frames that are interchanged every now and then.
	// 
	// NewAnimation() accepts 3 parameters:
	// ganim8.NewAnimation(frames, durations, onLoop)
	// 
	// frames is an array of frames ([]*image.Rectangle). You could provide
	// your own slice if you wanted to, but using a grid to get them is very convenient.
	// 
	// durations is a time.Duration or []time.Duration or
	// map[string]time.Duration.
	// When it's a time.Duration, it represents the duration of all frames
	// in the animation.
	// When it's a slice, it can represent different durations for different frames.
	// You can specify durations for all frames individually or you can
	// specify duration for ranges of frames:
	// map[string]time.Duration{ "3-5" : 2 * time.Millisecond }
	// 
	// onLoop is function of the animetion methods or callback.
	// If ganim8.Nop is specified, it does nothing.
	// If ganim8.PauseAtEnd is specified, it pauses at the end of the animation.
	// It can be any function that follows the type "func(anim *Animation, loops int)".
	// The first parameter is the animation object and the second parameter is
	// the count of the loops that elapsed since the previous Animation.Update().
	g.anim = ganim8.NewAnimation(frames, 100*time.Millisecond, ganim8.Nop)

	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	// Animation.Draw() draws the current frame of the animation.
	// It accepts screen image to draw on, source texture image, and draw options.
	// Draw options are x, y, angle (radian), scaleX, scaleY, originX and originY.
	// OriginX and OriginY are useful to draw the animation with scaling, centering,
	// rotating etc.
	// 
	// DrawOpts() is just a shortcut of creating DrawOption object.
	// It only needs the first 2 parameters x and y.
	// The rest of the parameters (angle, scaleX, scaleY, originX, orignY)
	// are optional.
	// If those are not specified, defaults values will be applied.
	// 
	// In this example, it draws the animation at the center of the screen.
	g.anim.Draw(screen, monsterImage, ganim8.DrawOpts(screenWidth/2, screenHeight/2, 0, 1, 1, 0.5, 0.5))
}

func (g *Game) Update() error {
	now := time.Now()

	// Animation.Update() updates the animation.
	// It receives time.Duration value and set the current frame of the animation. 
	// Each duration time of each frames can be customized for example like this:
	// ganim8.NewAnimation(
	//   grid.GetFrames("1-5", 5), 
	//   map[string]time.Duration {
	//     "1-2" : 100*time.Millisecond,
	//     "3"   : 300*time.Millisecond,
	//     "4-5" : 100*time.Millisecond,
	// })
	g.anim.Update(now.Sub(g.prevUpdateTime))

	g.prevUpdateTime = now
	return nil
}
```

[source code](https://github.com/yohamta/ganim8/blob/master/examples/simple/main.go)

### Output

<p align="center">
  <img src="https://github.com/yohamta/ganim8/blob/master/examples/gif/example.gif?raw=true" />
</p>

### The texture used in the example

  <img src="https://github.com/yohamta/ganim8/blob/master/examples/assets/images/Character_Monster_Slime_Blue.png?raw=true" />
