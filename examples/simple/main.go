package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/ganim8/v2"
	"github.com/yohamta/ganim8/v2/examples/assets/images"
)

const (
	screenWidth  = 300
	screenHeight = 300
)

type Game struct {
	prevUpdateTime time.Time
	anim           *ganim8.Animation
}

var (
	monsterImage = ebiten.NewImageFromImage(bytes2Image(&images.CHARACTER_MONSTER_SLIME_BLUE))
)

func (g *Game) Update() error {
	elapsedTime := time.Now().Sub(g.prevUpdateTime)

	g.anim.Update(elapsedTime)

	g.prevUpdateTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.anim.Draw(screen, ganim8.DrawOpts(screenWidth/2, screenHeight/2, 0, 1, 1, 0.5, 0.5))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{
		prevUpdateTime: time.Now(),
	}
	g.setupAnimations()

	return g
}

func (g *Game) setupAnimations() {
	grid := ganim8.NewGrid(100, 100, 500, 600)
	sprite := ganim8.NewSprite(monsterImage, grid.GetFrames("1-5", 5))
	g.anim = ganim8.NewAnimation(sprite, 100*time.Millisecond, ganim8.Nop)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func bytes2Image(rawImage *[]byte) image.Image {
	img, format, error := image.Decode(bytes.NewReader(*rawImage))
	if error != nil {
		log.Fatal("Bytes2Image Failed: ", format, error)
	}
	return img
}
