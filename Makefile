CSS      = reveal.js/css/...
JS       = reveal.js/js/...
LIB      = reveal.js/lib/...
PLUGIN   = reveal.js/plugin/...
THEMES	 = themes/...
ASSETS	 = assets/...

bindata:
	$(GOPATH)/bin/go-bindata -o=src/github.com/nautsio/presenter/assets.go $(CSS) $(JS) $(LIB) $(PLUGIN) $(THEMES) $(ASSETS)

release: osx linux windows

osx:
	GOOS=darwin GOARCH=amd64 gb build all

linux:
	GOOS=linux GOARCH=amd64 gb build

windows:
	GOOS=windows GOARCH=amd64 gb build
