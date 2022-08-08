package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/ganim8/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	prev          time.Time
	spinning      []*ganim8.Animation
	plane         *ganim8.Animation
	seaplane      *ganim8.Animation
	seaplaneAngle float64
	submarine     *ganim8.Animation
}

func NewGame() *Game {
	g := &Game{
		prev: time.Now(),
	}
	g.setupAnimations()

	return g
}

func (g *Game) Update() error {
	now := time.Now()
	delta := now.Sub(g.prev)
	g.prev = now

	for _, a := range g.spinning {
		a.Update(delta)
	}
	g.plane.Update(delta)
	g.seaplane.Update(delta)
	g.submarine.Update(delta)
	g.seaplaneAngle += g.seaplaneAngle + float64(delta.Milliseconds())*math.Pi/180

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	for i, a := range g.spinning {
		a.Draw(screen, ganim8.DrawOpts(float64(i)*75, float64(i)*50))
		// Alternative way to draw an animation:
		// ganim8.DrawAnime(screen, a, float64(i)*75, float64(i)*50, 0, 1, 1, .5, .5)
	}
	g.plane.Draw(screen, ganim8.DrawOpts(100, 400))
	g.seaplane.Draw(screen, ganim8.DrawOpts(250, 432, g.seaplaneAngle, 1, 1, 32, 32))
	g.submarine.Draw(screen, ganim8.DrawOpts(600, 100))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) setupAnimations() {
	img := ebiten.NewImageFromImage(readImage("assets/1945.png"))

	//                    frame(w,h), image(w,h), offsets, border
	g32 := ganim8.NewGrid(32, 32, 1024, 768, 3, 3, 1)

	g.spinning = []*ganim8.Animation{
		ganim8.New(img, g32.Frames("1-8", 1), time.Millisecond*100),
		ganim8.New(img, g32.Frames(18, "8-11", 18, "10-7"), time.Millisecond*200),
		ganim8.New(img, g32.Frames("1-8", 2), time.Millisecond*300),
		ganim8.New(img, g32.Frames(19, "8-11", 19, "10-7"), time.Millisecond*400),
		ganim8.New(img, g32.Frames("1-8", 3), time.Millisecond*500),
		ganim8.New(img, g32.Frames(20, "8-11", 20, "10-7"), time.Millisecond*600),
		ganim8.New(img, g32.Frames("1-8", 4), time.Millisecond*700),
		ganim8.New(img, g32.Frames(21, "8-11", 21, "10-7"), time.Millisecond*800),
		ganim8.New(img, g32.Frames("1-8", 5), time.Millisecond*900),
	}

	//                    frame(w,h), image(w,h), offsets, border
	g64 := ganim8.NewGrid(64, 64, 1024, 768, 299, 101, 2)

	g.plane = ganim8.New(img, g64.Frames(1, "1-3"), time.Millisecond*100)
	g.seaplane = ganim8.New(img, g64.Frames("2-4", 3), time.Millisecond*100)
	g.seaplaneAngle = 0

	//                   frame(w,h), image(w,h), offsets, border
	gs := ganim8.NewGrid(32, 98, 1024, 768, 366, 102, 1)

	g.submarine = ganim8.New(img, gs.Frames("1-7", 1, "6-2", 1),
		// individual frame delays
		map[string]time.Duration{
			"1":    time.Second * 2,
			"2-6":  time.Millisecond * 100,
			"7":    time.Second * 1,
			"8-12": time.Millisecond * 100,
		})
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
