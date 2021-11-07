package ganim8

import "github.com/hajimehoshi/ebiten/v2"

// DrawOptions represents the option for Sprite.Draw().
// For shortcut, DrawOpts() function can be used.
type DrawOptions struct {
	X, Y             float64
	Rotate           float64
	ScaleX, ScaleY   float64
	OriginX, OriginY float64
	ColorM           ebiten.ColorM
	CompositeMode    ebiten.CompositeMode
}

// ShaderOptions represents the option for Sprite.DrawWithShader()
type ShaderOptions struct {
	Uniforms map[string]interface{}
	Shader   *ebiten.Shader
	Images   [3]*ebiten.Image
}

// DrawOpts returns DrawOptions pointer with specified
// settings.
// The paramters are x, y, rotate (in radian), scaleX, scaleY
// originX, originY.
// If scaleX and ScaleY is not specified the default value
// will be 1.0, 1.0.
// If OriginX and OriginY is not specified the default value
// will be 0, 0
func DrawOpts(x, y float64, args ...float64) *DrawOptions {
	r, sx, sy, ox, oy := 0., 1., 1., 0., 0.
	switch len(args) {
	case 5:
		oy = args[4]
		fallthrough
	case 4:
		ox = args[3]
		fallthrough
	case 3:
		sy = args[2]
		fallthrough
	case 2:
		sx = args[1]
		fallthrough
	case 1:
		r = args[0]
	}
	return &DrawOptions{
		X:             x,
		Y:             y,
		Rotate:        r,
		ScaleX:        sx,
		ScaleY:        sy,
		OriginX:       ox,
		OriginY:       oy,
		ColorM:        ebiten.ColorM{},
		CompositeMode: ebiten.CompositeModeSourceOver,
	}
}
