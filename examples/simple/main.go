package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/ganim8/v2"
)

const (
	screenWidth  = 300
	screenHeight = 300
)

type Game struct {
	prev      time.Time
	animation *ganim8.Animation
}

func NewGame() *Game {
	g := &Game{
		prev: time.Now(),
	}
	g.setupMonsterAnimation()

	return g
}

func (g *Game) Update() error {
	elapsedTime := time.Now().Sub(g.prev)
	g.prev = time.Now()

	g.animation.Update(elapsedTime)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.animation.Draw(screen, ganim8.DrawOpts(screenWidth/2, screenHeight/2, 0, 1, 1, 0.5, 0.5))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) setupMonsterAnimation() {
	grid := ganim8.NewGrid(100, 100, 500, 600)
	img := ebiten.NewImageFromImage(readImage("assets/monster.png"))
	sprite := ganim8.NewSprite(img, grid.GetFrames("1-5", 5))
	g.animation = ganim8.NewAnimation(sprite, 100*time.Millisecond, ganim8.Nop)
}

//go:embed assets/*
var assets embed.FS

func readImage(file string) image.Image {
	b, _ := assets.ReadFile(file)
	return bytes2Image(&b)
}

func bytes2Image(rawImage *[]byte) image.Image {
	img, format, error := image.Decode(bytes.NewReader(*rawImage))
	if error != nil {
		log.Fatal("Bytes2Image Failed: ", format, error)
	}
	return img
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
