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
package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	resources "github.com/hajimehoshi/ebiten/examples/resources/images/shader"
	"github.com/tomlister/kager/editor"
	"github.com/tomlister/kager/viewer"
	"golang.org/x/image/font"
)

// Game implements the game interface.
type Game struct {
	Editor editor.Editor
	Viewer viewer.Viewer
}

// Update is called every 1/60th of a second
func (g *Game) Update(screen *ebiten.Image) error {
	g.Viewer.Time++
	g.Editor.Logic()
	return nil
}

// Draw is called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.Editor.Render(screen)
	g.Viewer.Render(g.Editor.Data, screen)
}

// Layout manages screen sizing
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 480
}

var (
	defaultFont   font.Face
	gopherImage   *ebiten.Image
	gopherBgImage *ebiten.Image
	normalImage   *ebiten.Image
	noiseImage    *ebiten.Image
)

func init() {
	// Import fonts
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 100
	defaultFont = truetype.NewFace(tt, &truetype.Options{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	// Import images
	img, _, err := image.Decode(bytes.NewReader(resources.Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	img, _, err = image.Decode(bytes.NewReader(resources.GopherBg_png))
	if err != nil {
		log.Fatal(err)
	}
	gopherBgImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	img, _, err = image.Decode(bytes.NewReader(resources.Normal_png))
	if err != nil {
		log.Fatal(err)
	}
	normalImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	img, _, err = image.Decode(bytes.NewReader(resources.Noise_png))
	if err != nil {
		log.Fatal(err)
	}
	noiseImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func main() {
	shaderEditor := editor.Editor{
		Data:  strings.Split(strings.ReplaceAll(string(radialblur_go), "\t", "        "), "\n"),
		Fonts: make([]*font.Face, 0),
	}
	shaderEditor.Fonts = append(shaderEditor.Fonts, &defaultFont)
	shaderViewer := viewer.Viewer{
		Images: make([]*ebiten.Image, 0),
	}
	shaderViewer.Fonts = append(shaderViewer.Fonts, &defaultFont)
	shaderViewer.Images = append(shaderViewer.Images, gopherImage)
	shaderViewer.Images = append(shaderViewer.Images, gopherBgImage)
	shaderViewer.Images = append(shaderViewer.Images, normalImage)
	shaderViewer.Images = append(shaderViewer.Images, noiseImage)
	game := &Game{
		Editor: shaderEditor,
		Viewer: shaderViewer,
	}

	ebiten.SetWindowSize(1280, 480)
	ebiten.SetWindowTitle("kager - shader editor")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
