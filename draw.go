package ganim8

import "github.com/hajimehoshi/ebiten/v2"

var drawOpts = DrawOpts(0, 0, 0, 1, 1, 1, 0.5, 0.5)

// DrawSprite draws a sprite to the screen.
func DrawSprite(screen *ebiten.Image, spr *Sprite, index int, x, y, rot, sx, sy, ox, oy float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rot)
	drawOpts.SetScale(sx, sy)
	drawOpts.SetOrigin(ox, oy)
	spr.Draw(screen, index, drawOpts)
}

// DrawSpriteWithOpts draws a sprite to the screen.
func DrawSpriteWithOpts(screen *ebiten.Image, spr *Sprite, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	if shaderOpts != nil {
		spr.DrawWithShader(screen, index, opts, shaderOpts)
	} else {
		spr.Draw(screen, index, opts)
	}
}

// DrawAnime draws an animation to the screen.
func DrawAnime(screen *ebiten.Image, anim *Animation, x, y, rot, sx, sy, ox, oy float64) {
	drawOpts.SetPos(x, y)
	drawOpts.SetRot(rot)
	drawOpts.SetScale(sx, sy)
	drawOpts.SetOrigin(ox, oy)
	anim.Draw(screen, drawOpts)
}

// DrawAnimeWithOpts draws an anime to the screen.
func DrawAnimeWithOpts(screen *ebiten.Image, anim *Animation, index int, opts *DrawOptions, shaderOpts *ShaderOptions) {
	if shaderOpts != nil {
		anim.DrawWithShader(screen, opts, shaderOpts)
	} else {
		anim.Draw(screen, opts)
	}
}
