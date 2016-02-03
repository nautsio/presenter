package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

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
		cli.StringFlag{
			Name:  "theme, t",
			Usage: "Use one of the built in themes",
		},
	},
}

func doInit(c *cli.Context) {
	presentationPath, _ := os.Getwd()
	if len(c.Args()) > 0 {
		presentationPath = c.Args()[0]
	}

	createExamplePresentationPath(presentationPath)
	createExampleSlides(presentationPath)
	createExampleTheme(presentationPath)
	createExampleImageDirectory(presentationPath)
}

func createExamplePresentationPath(presentationPath string) {
	error := os.Mkdir(presentationPath, 0777)
	if error != nil {
		log.Fatal(error)
	}
}

func createExampleSlides(presentationPath string) {
	// Create slides file.
	file, error := os.Create(path.Join(presentationPath, "slides.md"))
	if error != nil {
		log.Fatal(error)
	}

	// Write contents of example slides to slides file.
	writer := bufio.NewWriter(file)
	slides, _ := Asset("assets/slides.md")
	writer.Write(slides)
	writer.Flush()
}

func createExampleImageDirectory(presentationPath string) {
	error := os.Mkdir(path.Join(presentationPath, "img"), 0777)
	if error != nil {
		log.Fatal(error)
	}
}

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
	theme, _ := Asset("assets/theme.css")
	writer.Write(theme)
	writer.Flush()
}

func doServe(c *cli.Context) {
	master := c.Bool("master")
	theme := c.String("theme")

	// Get the presentation path from the command line, or grab the current directory.
	presentationPath, _ := os.Getwd()
	if len(c.Args()) > 0 {
		presentationPath = c.Args()[0]
	}

	// Check if the path is relative.
	if !strings.HasPrefix(presentationPath, "/") {
		presentationPath, _ = filepath.Abs(presentationPath)
	}

	// Check if there is a presentation file present.
	if _, err := os.Stat(path.Join(presentationPath, "slides.md")); err != nil {
		fmt.Printf("slides.md does not exist at %s\n", presentationPath)
		os.Exit(1)
	}

	if theme != "" {
		fmt.Println("Using one of the packaged themes ...")
		http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, "themes/" + theme}))
		http.Handle("/fonts/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, "themes/" + theme}))
	} else {
		if _, err := os.Stat(path.Join(presentationPath, "css", "theme.css")); err == nil {
			fmt.Println("Found a theme, using it ...")
			http.Handle("/css/", http.FileServer(http.Dir(presentationPath)))
			http.Handle("/fonts/", http.FileServer(http.Dir(presentationPath)))
		} else {
			fmt.Println("No theme found ...")
		}
	}

	// Handle the slides.
	http.HandleFunc("/slides.md", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles(path.Join(presentationPath, "slides.md"))
		t.Execute(w, "")
	})

	// Handle images.
	http.HandleFunc("/img/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(presentationPath, r.URL.Path))
	})

	// Handle reveal.js files.
	http.Handle("/reveal.js/css/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, ""}))
	http.Handle("/reveal.js/js/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, ""}))
	http.Handle("/reveal.js/lib/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, ""}))
	http.Handle("/reveal.js/plugin/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, ""}))

	// Handle the website.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opt := &Option{Markdown: "slides.md", Master: master}
		indexHTML, _ := Asset("assets/index.html")
		indexTemplate := template.Must(template.New("index").Parse(string(indexHTML)))
		indexTemplate.Execute(w, opt)
	})

	fmt.Println("Opening browser and redirecting to the presentation ...")
	browser.OpenURL("http://localhost:8989")

	err := http.ListenAndServe(":8989", nil)
	panic("Error while serving slides: " + err.Error())
}
