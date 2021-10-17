# ganim8

Animation library for [Ebiten](https://ebiten.org/) which is a golang version of [anim8](https://github.com/kikito/anim8).

[GoDoc](https://pkg.go.dev/github.com/miyahoyo/ganim8)

## Simple Example

### Texture
  <img src="https://github.com/miyahoyo/ganim8/blob/master/examples/assets/images/Character_Monster_Slime_Blue.png?raw=true" />

### Source
```go
type Game struct {
	prevUpdateTime time.Time
	animes         []*ganim8.Animation
}

func NewGame() *Game {
	g := &Game{
		prevUpdateTime: time.Now(),
	}
	g.setupAnimations()

	return g
}

func (g *Game) setupAnimations() {
	// Make a new grid with the frame size of 100x100 and 
	// the image size of 500x600
	grid := ganim8.NewGrid(100, 100, 500, 600)
	// Make an animation from the grid
	// This specifies 1-5 columns, 5 row
	g.anim = ganim8.NewAnimation(grid.GetFrames("1-5", 5), 100*time.Millisecond, ganim8.Nop)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	// Draw the animation at the center of the screen
	g.anim.Draw(screen, monsterImage, ganim8.DrawOpts(screenWidth/2, screenHeight/2, 0, 1, 1, 0.5, 0.5))
}

func (g *Game) Update() error {
	elapsedTime := time.Now().Sub(g.prevUpdateTime)

	// Update the animation
	g.anim.Update(elapsedTime)

	g.prevUpdateTime = time.Now()
	return nil
}
```

[source code](https://github.com/miyahoyo/ganim8/blob/master/examples/simple/main.go)

### Output

<p align="center">
  <img src="https://github.com/miyahoyo/ganim8/blob/master/examples/gif/example.gif?raw=true" />
</p>
