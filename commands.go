package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"

	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/pkg/browser"
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
	Usage:       "",
	Description: "",
	Action:      doInit,
}

var commandServe = cli.Command{
	Name:        "serve",
	Usage:       "",
	Description: "",
	Action:      doServe,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "master, m",
			Usage: "Start presenter in master mode",
		},
	},
}

func doInit(c *cli.Context) {
	// Create slides file.
	file, error := os.Create("slides.md")
	if error != nil {
		log.Fatal(error)
	}

	// Write contents of example slides to slides file.
	writer := bufio.NewWriter(file)
	slides, _ := Asset("theme/slides.md")
	writer.Write(slides)
	writer.Flush()
}

func doServe(c *cli.Context) {
	slidesPath := c.Args()[0]
	slidesFile := filepath.Base(slidesPath)
	slidesWD, _ := os.Getwd()

	master := c.Bool("master")

	// Handle the slides
	http.HandleFunc("/"+slidesFile,
		func(w http.ResponseWriter, r *http.Request) {
			t, _ := template.ParseFiles(slidesPath)
			t.Execute(w, "")
		})

	// Handle the images
	http.HandleFunc("/img/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(slidesWD, r.URL.Path))
	})

	http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/js/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/lib/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/plugin/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, "reveal.js"}))
	http.Handle("/theme/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opt := &Option{Markdown: slidesFile, Master: master}
		indexHTML, _ := Asset("theme/index.html")
		indexTemplate := template.Must(template.New("index").Parse(string(indexHTML)))
		indexTemplate.Execute(w, opt)
	})

	fmt.Println("Opening browser and redirecting to the presentation ...")
	browser.OpenURL("http://localhost:8989")

	err := http.ListenAndServe(":8989", nil)
	panic("Error while serving slides: " + err.Error())
}
