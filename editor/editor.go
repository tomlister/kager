/*
	Copyright 2020 Tom Lister & Kager Authors

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package editor

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// Vec2 contains x and y of a point in space or 2d array
type Vec2 struct {
	x, y float64
}

// Rect contains x, y, width and height values
type Rect struct {
	x, y, w, h float64
}

// Editor contains the shader data, fonts and other information
type Editor struct {
	Data           []string
	Fonts          []*font.Face
	ScrollOffset   Vec2
	CursorPos      Vec2
	KeyInterval    int
	CursorInterval int
}

func (r Rect) collideVec2(p Vec2) bool {
	if r.x < p.x &&
		r.x+r.w > p.x &&
		r.y < p.y &&
		r.y+r.h > p.y {
		return true
	}
	return false
}

func (e *Editor) findCursorPos(point Vec2) (Vec2, bool) {
	for i, data := range e.Data {
		// TODO: Figure a less shitty way of calculating the bounds of dataless lines
		if data == "" {
			data = "I"
		}
		bounds := text.BoundString((*e.Fonts[0]), data)
		rect := Rect{20, 20 + e.ScrollOffset.y + float64(i*bounds.Dy()), 640, float64(bounds.Dy())}
		if rect.collideVec2(point) {
			return Vec2{0, float64(i)}, true
		}
	}
	return Vec2{}, false
}

func (e *Editor) addWhiteSpace() {
	for i, data := range e.Data {
		if !strings.HasSuffix(data, " ") {
			e.Data[i] += " "
		}
	}
}

// Logic handles editor logic
func (e *Editor) Logic() {
	e.addWhiteSpace()
	_, yoffset := ebiten.Wheel()
	e.ScrollOffset.y += yoffset
	// TODO: Fix clicking
	/*if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		pos, found := e.findCursorPos(Vec2{float64(cx), float64(cy)})
		if found {
			e.CursorPos = pos
		}
	}*/
	if e.KeyInterval == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			e.KeyInterval = 10
			if len(e.Data)-1 > int(e.CursorPos.y)+1 {
				if int(e.CursorPos.x) > len(e.Data[int(e.CursorPos.y)+1]) {
					e.CursorPos.x = float64(len(e.Data[int(e.CursorPos.y)+1])) - 1
				}
				e.CursorPos.y++
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
			e.KeyInterval = 10
			if int(e.CursorPos.x) > len(e.Data[int(e.CursorPos.y)-1]) {
				e.CursorPos.x = float64(len(e.Data[int(e.CursorPos.y)-1])) - 1
			}
			e.CursorPos.y--
		} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			e.KeyInterval = 10
			if e.CursorPos.x > 0 {
				e.CursorPos.x--
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
			e.KeyInterval = 10
			if int(e.CursorPos.x) < len(e.Data[int(e.CursorPos.y)])-1 && len(e.Data[int(e.CursorPos.y)]) > 0 {
				e.CursorPos.x++
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			e.KeyInterval = 5
			if e.CursorPos.x > 0 {
				line := e.Data[int(e.CursorPos.y)]
				e.Data[int(e.CursorPos.y)] = line[:int(e.CursorPos.x)-1] + line[int(e.CursorPos.x):]
				e.CursorPos.x--
			} else {
				copy(e.Data[int(e.CursorPos.y):], e.Data[int(e.CursorPos.y)+1:])
				e.Data[len(e.Data)-1] = ""
				e.Data = e.Data[:len(e.Data)-1]
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			e.KeyInterval = 5
			e.Data = append(e.Data[:int(e.CursorPos.y)+2], e.Data[int(e.CursorPos.y+1):]...)
			e.Data[int(e.CursorPos.y)+1] = ""
		}
	} else {
		e.KeyInterval--
	}
	line := e.Data[int(e.CursorPos.y)]
	inchars := string(ebiten.InputChars())
	e.Data[int(e.CursorPos.y)] = line[:int(e.CursorPos.x)] + inchars + line[int(e.CursorPos.x):]
	e.CursorPos.x += float64(len(inchars))
}

// Render draws the editor
func (e *Editor) Render(screen *ebiten.Image) {
	bounds := text.BoundString((*e.Fonts[0]), e.Data[0])
	// TODO: Make Cursor blinking elegant
	if e.CursorInterval > 0 {
		e.CursorInterval--
	} else if e.CursorInterval == 0 {
		e.CursorInterval = 30
	}
	for i, data := range e.Data {
		if int(e.CursorPos.y) == i {
			opts := &ebiten.DrawImageOptions{}
			activeBackground, _ := ebiten.NewImage(640, bounds.Dy(), ebiten.FilterDefault)
			activeBackground.Fill(color.RGBA{0x32, 0x32, 0x32, 0xff})
			opts.GeoM.Translate(0, 5+(e.ScrollOffset.y)+float64(i*bounds.Dy()))
			screen.DrawImage(activeBackground, opts)
			activeBackground.Dispose()
			if e.CursorInterval > 15 {
				if len(data) == 0 {
					data += "|"
				} else {
					runeData := []rune(data)
					runeData[int(e.CursorPos.x)] = rune('|')
					data = string(runeData)
				}
			}
			text.Draw(screen, data, (*e.Fonts[0]), 20, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), color.White)
		} else {
			text.Draw(screen, data, (*e.Fonts[0]), 20, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), color.White)
		}
		text.Draw(screen, fmt.Sprintf("%d", 1+i), (*e.Fonts[0]), 2, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), color.White)
	}
}
