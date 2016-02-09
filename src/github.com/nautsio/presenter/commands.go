package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
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

// Built in themes.
var themes = []string{"nauts", "xebia"}

// Commands is an array containing the available commands.
var Commands = []cli.Command{
	commandInit,
	commandServe,
}

var commandInit = cli.Command{
	Name:      "init",
	ShortName: "i",
	Usage:     "Initialize an empty presentation directory",
	ArgsUsage: "<destination path>",
	Action:    doInit,
}

var commandServe = cli.Command{
	Name:      "serve",
	ShortName: "s",
	Usage:     "Serve a presentation directory",
	ArgsUsage: "<presentation directory>",
	Action:    doServe,
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

// Initialize an example presentation.
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
	slides, _ := Asset("assets/slides.md")
	writer.Write(slides)
	writer.Flush()
}

// Create the example presentation image directory.
func createExampleImageDirectory(presentationPath string) {
	error := os.Mkdir(path.Join(presentationPath, "img"), 0777)
	if error != nil {
		log.Fatal(error)
	}
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
	theme, _ := Asset("assets/theme.css")
	writer.Write(theme)
	writer.Flush()
}

// Check if the given theme exists.
func themeExists(theme string) bool {
	sort.Strings(themes)
	i := sort.SearchStrings(themes, theme)
	return (i < len(themes) && themes[i] == theme)
}

// Check if the given file exists.
func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err != nil
}

// Check if the given path is absolute.
func pathIsAbsolute(path string) bool {
	return strings.HasPrefix(path, "/")
}

func concatFiles(path string, extension string) []byte {
	d, _ := os.Open(path)
	defer d.Close()
	files, _ := d.Readdir(-1)
	var content []byte
	for _, file := range files {
		if file.Mode().IsRegular() && filepath.Ext(file.Name()) == extension {
			c, _ := ioutil.ReadFile(file.Name())
			content = append(content, '\n')
			content = append(content, c...)
			content = append(content, '\n')
		}
	}
	return content
}

// Serve the presentation.
func doServe(c *cli.Context) {
	master := c.Bool("master")
	theme := c.String("theme")

	// Get the presentation path from the command line, or grab the current directory.
	presentationPath, _ := os.Getwd()
	if len(c.Args()) > 0 {
		presentationPath = c.Args()[0]
	}

	// Check if the path is relative.
	if !pathIsAbsolute(presentationPath) {
		presentationPath, _ = filepath.Abs(presentationPath)
	}

	// Check if there is a presentation file present.
	// if fileExists(path.Join(presentationPath, "slides.md")) {
	// 	fmt.Printf("slides.md does not exist at %s\n", presentationPath)
	// 	os.Exit(1)
	// }

	// Check if a theme was passed.
	if themeExists(theme) {
		fmt.Println("Using one of the packaged themes ...")
		http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, "themes/" + theme}))
		http.Handle("/fonts/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, AssetInfo, "themes/" + theme}))
	} else {
		if !fileExists(path.Join(presentationPath, "css", "theme.css")) {
			fmt.Println("Found a theme, using it ...")
			http.Handle("/css/", http.FileServer(http.Dir(presentationPath)))
			http.Handle("/fonts/", http.FileServer(http.Dir(presentationPath)))
		} else {
			fmt.Println("No theme found ...")
		}
	}

	// Handle the slides.
	http.HandleFunc("/slides.md", func(w http.ResponseWriter, r *http.Request) {
		slides := concatFiles(presentationPath, ".md")
		fmt.Printf("%s\n", slides)
		w.Write(slides)
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
