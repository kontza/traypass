cmd/list_templ.go: cmd/list.templ
	templ generate cmd/list.templ

.PHONY=build
build: build/bin/waitwhat.app

build/bin/waitwhat.app: main.go cmd/*.go
	wails build -ldflags '-w -s'
