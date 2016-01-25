# Presenter
The idea is to have a single binary that can present your reveal.js slides with the nauts.io theme (or custom, if you edit the theme.css file). 

We also want it to be a hosted service where you can add a markdown file hosted on github (or somewhere else) and it will serve the slides for you online. (e.g. slides.nauts.io/github.com/nautsio/presenter/templates/slides.md).

To generate the assets.go:
```
make
```

Usage:
```
// Generate slides.md
presenter init

// Start the internal webserver
presenter serve slides.md
```
