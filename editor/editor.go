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
				if len(e.Data[int(e.CursorPos.y)]) > 0 {
					if e.CursorPos.y > 0 {
						e.Data[int(e.CursorPos.y-1)] += e.Data[int(e.CursorPos.y)]
					}
				}
				copy(e.Data[int(e.CursorPos.y):], e.Data[int(e.CursorPos.y)+1:])
				e.Data[len(e.Data)-1] = ""
				e.Data = e.Data[:len(e.Data)-1]
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			e.KeyInterval = 5
			e.Data = append(e.Data[:int(e.CursorPos.y)+2], e.Data[int(e.CursorPos.y+1):]...)
			line := e.Data[int(e.CursorPos.y)]
			e.Data[int(e.CursorPos.y)+1] = line[int(e.CursorPos.x):]
			e.Data[int(e.CursorPos.y)] = line[:int(e.CursorPos.x)]
			e.CursorPos.y++
			e.CursorPos.x = 0
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
	screen.Fill(color.RGBA{0x29, 0x2D, 0x3E, 0xff})
	for i, data := range e.Data {
		if int(e.CursorPos.y) == i {
			opts := &ebiten.DrawImageOptions{}
			activeBackground, _ := ebiten.NewImage(640, bounds.Dy(), ebiten.FilterDefault)
			activeBackground.Fill(color.RGBA{0x41, 0x48, 0x63, 0xff})
			opts.GeoM.Translate(0, 5+(e.ScrollOffset.y)+float64(i*bounds.Dy()))
			screen.DrawImage(activeBackground, opts)
			activeBackground.Dispose()
			syntax := hightLightLine(data)
			syntax.drawHighLighted(e, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), screen)
			if e.CursorInterval > 15 {
				xoffset := 0
				xprev := 0
				for i, character := range data {
					if i == int(e.CursorPos.x+1) {
						break
					}
					bounds := text.BoundString((*e.Fonts[0]), string(character))
					fmt.Printf("%s, %d\n", string(character), bounds.Dx())
					if bounds.Dx() > 0 {
						xoffset += bounds.Dx()
						xprev = bounds.Dx()
					} else {
						xoffset += 4
						xprev = 4
					}
				}
				opts = &ebiten.DrawImageOptions{}
				cursor, _ := ebiten.NewImage(1, bounds.Dy(), ebiten.FilterDefault)
				cursor.Fill(color.RGBA{0xab, 0x47, 0xbc, 0xff})
				opts.GeoM.Translate(float64(20+xoffset-xprev), 5+float64(int(e.ScrollOffset.y)+(i*int(bounds.Dy()))))
				screen.DrawImage(cursor, opts)
				cursor.Dispose()
			}
		} else {
			syntax := hightLightLine(data)
			syntax.drawHighLighted(e, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), screen)
		}
		text.Draw(screen, fmt.Sprintf("%d", 1+i), (*e.Fonts[0]), 2, 20+int(e.ScrollOffset.y)+(i*int(bounds.Dy())), color.White)
	}
}
