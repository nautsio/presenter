package command

import (
	"bufio"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/nautsio/presenter/bindata"
)

// Init creates an example presentation.
func Init(c *cli.Context) {
	presentationPath, _ := os.Getwd()
	if len(c.Args()) > 0 {
		presentationPath = c.Args()[0]
	}

	createExamplePresentationPath(presentationPath)
	createExampleSlides(presentationPath)
	createExampleTheme(presentationPath)
	createExampleImageDirectory(presentationPath)
}

// Create the example presentation directory.
func createExamplePresentationPath(presentationPath string) {
	error := os.Mkdir(presentationPath, 0777)
	if error != nil {
		log.Fatal(error)
	}
}

// Create the example presentation slides.
func createExampleSlides(presentationPath string) {
	// Create slides file.
	file, error := os.Create(path.Join(presentationPath, "slides.md"))
	if error != nil {
		log.Fatal(error)
	}

	// Write contents of example slides to slides file.
	writer := bufio.NewWriter(file)
	slides, _ := bindata.Asset("assets/slides.md")
	writer.Write(slides)
	writer.Flush()
}

// Create the example presentation theme.
func createExampleTheme(presentationPath string) {
	error := os.Mkdir(path.Join(presentationPath, "css"), 0777)
	if error != nil {
		log.Fatal(error)
	}

	// Create theme file.
	file, error := os.Create(path.Join(presentationPath, "css", "theme.css"))
	if error != nil {
		log.Fatal(error)
	}

	// Write contents of example theme to theme file.
	writer := bufio.NewWriter(file)
	theme, _ := bindata.Asset("assets/theme.css")
	writer.Write(theme)
	writer.Flush()
}

// Create the example presentation image directory.
func createExampleImageDirectory(presentationPath string) {
	error := os.Mkdir(path.Join(presentationPath, "img"), 0777)
	if error != nil {
		log.Fatal(error)
	}
}
