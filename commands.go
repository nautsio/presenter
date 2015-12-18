package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/elazarl/go-bindata-assetfs"
)

// Option is passed to the templating engine.
type Option struct {
	Markdown string
	Master   bool
}

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
	mdFilePath := c.Args()[0]
	mdFileName := filepath.Base(mdFilePath)

	// markdown file
	http.HandleFunc("/"+mdFileName,
		func(w http.ResponseWriter, r *http.Request) {
			t, _ := template.ParseFiles(mdFilePath)
			t.Execute(w, "")
		})

	http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/js/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/lib/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/plugin/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/assets/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opt := &Option{Markdown: mdFileName, Master: false}
		indexHTML, _ := Asset("templates/index.html")
		indexTemplate := template.Must(template.New("index").Parse(string(indexHTML)))
		indexTemplate.Execute(w, opt)
	})

	fmt.Println("Presenting at http://localhost:8989")

	err := http.ListenAndServe(":8989", nil)
	panic("Error while serving slides: " + err.Error())
}
