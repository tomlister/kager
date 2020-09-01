package editor

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

// Button stores button information
type Button struct {
	Position   Vec2
	Text       string
	Color      color.RGBA
	ColorHover color.RGBA
	Callback   func(editor *Editor)
}

// Test handles clicking of the button
func (b Button) Test(editor *Editor) {
	bounds := text.BoundString((*editor.Fonts[0]), b.Text)
	rect := Rect{b.Position.X, b.Position.Y, float64(bounds.Dx()) + 10, float64(bounds.Dy()) + 10}
	mx, my := ebiten.CursorPosition()
	if rect.collideVec2(Vec2{float64(mx), float64(my)}) {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			b.Callback(editor)
		}
	}
}

// Draw renders the button
func (b Button) Draw(editor *Editor, screen *ebiten.Image) {
	bounds := text.BoundString((*editor.Fonts[0]), b.Text)
	rect := Rect{b.Position.X, b.Position.Y, float64(bounds.Dx()) + 10, float64(bounds.Dy()) + 10}
	mx, my := ebiten.CursorPosition()
	col := b.Color
	if rect.collideVec2(Vec2{float64(mx), float64(my)}) {
		col = b.ColorHover
	}
	background, _ := ebiten.NewImage(int(rect.w), int(rect.h), ebiten.FilterDefault)
	background.Fill(col)
	ops := &ebiten.DrawImageOptions{}
	ops.GeoM.Translate(rect.x, rect.y)
	screen.DrawImage(background, ops)
	background.Dispose()
	text.Draw(screen, b.Text, (*editor.Fonts[0]), int(rect.x)+5, int(rect.y)+bounds.Dy()+2, color.White)
}
