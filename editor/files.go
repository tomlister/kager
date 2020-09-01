package editor

import (
	"log"

	"github.com/sqweek/dialog"
)

// OpenFileButtonCallback creates a dialog when the button is pressed
func OpenFileButtonCallback(editor *Editor) {
	_, err := dialog.File().Filter("Kage Shader File (.go)", "go").Load()
	if err != nil {
		log.Fatalln(err)
	}
}
