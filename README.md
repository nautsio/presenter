# Presenter
If you have edited any of the theme or asset files, you will need to generate the assets.go file again.
```
// If you have updated reveal.js.
git submodule init
git submodule update

// Generate binary data.
make bindata
```

To build the project use gb.
```
go get github.com/constabulary/gb/...

gb vendor restore
gb build all
```

## Presentation structure
The structure of a presentation folder is expected to look similar to this:
```
├── css             (custom styling should go in theme.css)
│   └── theme.css
├── img             (images that are used in the presentation or self-made theme go here)
│   ├── logo.png
└── slides.md       (slides are expected to be in the slides.md file)
```

## Creating a presentation
**NOT YET WORKING CORRECTLY**   
An example presentation directory can be created by running the init command.
```
// Create a presentation directory called "presentation":
presenter init presentation

// or
presenter init /absolute/path/to/presentation

// or
presenter init relative/path/to/presentation
```


## Presenting a presentation
To view a presentation, point presenter at a directory containing your presentation.
```
// With an absolute path:
presenter serve /absolute/path/to/presentation

// With a relative path:
presenter serve relative/path/to/presentation

// Or by serving from the current directory:
presenter serve
```

## Themes
If you want to use one of the built in themes, supply the theme flag.
```
presenter serve -t nauts

// or
presenter serve --theme nauts
```
