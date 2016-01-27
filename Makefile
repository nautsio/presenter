CSS      = reveal.js/css/...
JS       = reveal.js/js/...
LIB      = reveal.js/lib/...
PLUGIN   = reveal.js/plugin/...
THEME		 = theme/...

bindata: assets.go
	$(GOPATH)/bin/go-bindata -o=assets.go $(CSS) $(JS) $(LIB) $(PLUGIN) $(THEME)

osx:
	go build -o $(GOPATH)/bin/presenter.osx .

linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(GOPATH)/bin/presenter.linux .

windows:
	GOOS=windows GOARCH=386 go build -o $(GOPATH)/bin/presenter.exe .
