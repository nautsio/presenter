package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"path/filepath"

	"github.com/codegangsta/cli"
)

// Commands is an array containing the available commands.
var Commands = []cli.Command{
	commandInit,
	commandServe,
}

var commandInit = cli.Command{
	Name:        "init",
	Usage:       "Create an example presentation",
	Description: "",
	Action:      doInit,
}

var commandServe = cli.Command{
	Name:        "serve",
	Usage:       "",
	Description: "",
	Action:      doServe,
}

func createDirectory(path string) {
	error := os.MkdirAll(path, 0755)
	if error != nil {
		log.Fatal(error)
	}
}

func createTheme(path string) {
	// Create theme styling.
	file, error := os.Create(path)
	if error != nil {
		log.Fatal(error)
	}

	// TODO: add the templates/slides.css into assets.
	// Write contents of example theme to theme file.
	writer := bufio.NewWriter(file)
	css, _ := Asset("templates/slides.css")
	writer.Write(css)
	writer.Flush()
}

func createPresentation(path string) {
	// Create slides file.
	file, error := os.Create(path)
	if error != nil {
		log.Fatal(error)
	}

	// TODO: add the templates/slides.md into assets.
	// Write contents of example slides to slides file.
	writer := bufio.NewWriter(file)
	slides, _ := Asset("templates/slides.md")
	writer.Write(slides)
	writer.Flush()
}

func doInit(c *cli.Context) {
	name := "presenter"
	if len(c.Args()) > 0 {
		name = c.Args()[0]
	}

	imgPath := filepath.Join(name, "img")
	cssPath := filepath.Join(name, "css")
	themePath := filepath.Join(cssPath, "theme.css")
	presentationPath := filepath.Join(name, "slides.md")

	createDirectory(cssPath)
	createDirectory(imgPath)
	createTheme(themePath)
	createPresentation(presentationPath)
}

func doServe(c *cli.Context) {
	fmt.Println("Presenting at http://localhost:8989")

	err := http.ListenAndServe(":8989", nil)
	panic("Error while serving slides: " + err.Error())
}
