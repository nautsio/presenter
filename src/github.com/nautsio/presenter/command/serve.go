package command

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/nautsio/presenter/bindata"
	"github.com/pkg/browser"
)

// Option is passed to the templating engine.
type Option struct {
	Markdown string
	Master   bool
}

// Built in themes.
var themes = []string{"nauts", "xebia"}

// Serve the presentation.
func Serve(c *cli.Context) {
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
	if fileExists(path.Join(presentationPath, "slides.md")) {
		log.Printf("slides.md does not exist at %s\n", presentationPath)
		os.Exit(1)
	}

	// Check if a theme was passed.
	if themeExists(theme) {
		log.Printf("Using build-in theme [%s]\n", theme)
		http.Handle("/css/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, "themes/" + theme}))
		http.Handle("/fonts/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, "themes/" + theme}))
	} else {
		if !fileExists(path.Join(presentationPath, "css", "theme.css")) {
			log.Println("Found a theme in the css directory, using it...")
			http.Handle("/css/", http.FileServer(http.Dir(presentationPath)))
			http.Handle("/fonts/", http.FileServer(http.Dir(presentationPath)))
		} else {
			log.Println("No theme found ...")
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
	http.Handle("/reveal.js/css/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, ""}))
	http.Handle("/reveal.js/js/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, ""}))
	http.Handle("/reveal.js/lib/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, ""}))
	http.Handle("/reveal.js/plugin/", http.FileServer(&assetfs.AssetFS{bindata.Asset, bindata.AssetDir, bindata.AssetInfo, ""}))

	// Handle the website.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opt := &Option{Markdown: "slides.md", Master: master}
		indexHTML, _ := bindata.Asset("assets/index.html")
		indexTemplate := template.Must(template.New("index").Parse(string(indexHTML)))
		indexTemplate.Execute(w, opt)
	})

	log.Println("Serving presentation on http://localhost:8989")
	log.Println("Opening browser and redirecting to the presentation, press Ctrl-C to stop ...")
	browser.OpenURL("http://localhost:8989")

	err := http.ListenAndServe(":8989", nil)
	log.Fatalf("Error while serving slides: %s\n", err.Error())
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
