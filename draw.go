package ganim8

import "github.com/hajimehoshi/ebiten/v2"

var drawOpts = DrawOpts(0, 0, 0, 1, 1, 1, 0.5, 0.5)

// DrawSpr draws a sprite to the screen.
func DrawSpr(screen *ebiten.Image, spr *Sprite, index int, x, y, rot, sx, sy, ox, oy float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rot)
	drawOpts.SetScale(sx, sy)
	drawOpts.SetOrigin(ox, oy)
	spr.Draw(screen, index, drawOpts)
}

// DrawAnime draws an animation to the screen.
func DrawAnime(screen *ebiten.Image, anim *Animation, x, y, rot, sx, sy, ox, oy float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rot)
	drawOpts.SetScale(sx, sy)
	drawOpts.SetOrigin(ox, oy)
	anim.Draw(screen, drawOpts)
}
