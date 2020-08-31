package viewer

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// Viewer holds information required to draw
type Viewer struct {
	Images []*ebiten.Image
	Fonts  []*font.Face
	Time   float32
}

// Render draws the shader and any errors
func (v *Viewer) Render(data []string, screen *ebiten.Image) {
	s, err := ebiten.NewShader([]byte(strings.Join(data[:], "\n")))
	if err != nil {
		text.Draw(screen, err.Error(), (*v.Fonts[0]), 640, 20, color.RGBA{0xff, 0x00, 0x00, 0xff})
	} else {
		w, h := 640, 480
		cx, cy := ebiten.CursorPosition()
		op := &ebiten.DrawRectShaderOptions{}
		op.GeoM.Translate(float64(640), float64(0))
		op.Uniforms = []interface{}{
			float32(v.Time) / 60,                // Time
			[]float32{float32(cx), float32(cy)}, // Cursor
			[]float32{float32(w), float32(h)},   // ScreenSize
		}
		op.Images[0] = v.Images[0] // gopherImage
		op.Images[1] = v.Images[2] // gopherBgImage
		op.Images[2] = v.Images[1] // normalImage
		op.Images[3] = v.Images[3] // noiseImage
		screen.DrawRectShader(w, h, s, op)
	}

}
