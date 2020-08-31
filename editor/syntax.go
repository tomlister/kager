package editor

import (
	"image/color"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

// LineSegment stores syntax highlighting information
type LineSegment struct {
	Data  string
	Color color.Color
}

// LineSegments is a slice of line segments
type LineSegments []LineSegment

func hightLightLine(line string) LineSegments {
	segments := make([]LineSegment, 0)
	//match uniforms
	matchUniforms, _ := regexp.MatchString("(var [a-zA-Z0-9_.-]+ [a-zA-Z0-9_.-]+)", line)
	if matchUniforms {
		splitstrings := strings.Split(line, " ")
		for i, split := range splitstrings {
			textcolor := color.RGBA{0xff, 0xff, 0xff, 0xff}
			switch i {
			case 0:
				textcolor = color.RGBA{0xc7, 0x92, 0xea, 0xff}
			case 1:
				textcolor = color.RGBA{0xff, 0xcb, 0x6b, 0xff}
			}
			segment := LineSegment{
				Data:  split,
				Color: textcolor,
			}
			segments = append(segments, segment)
		}
	}
	if len(segments) == 0 {
		segments = append(segments, LineSegment{
			Data:  line,
			Color: color.White,
		})
	}
	return segments
}

func (l LineSegments) drawHighLighted(e *Editor, y int, screen *ebiten.Image) {
	xoffset := 0
	for _, segment := range l {
		bounds := text.BoundString((*e.Fonts[0]), segment.Data)
		text.Draw(screen, segment.Data, (*e.Fonts[0]), 20+xoffset, y, segment.Color)
		xoffset += bounds.Dx() + 4
	}
}
